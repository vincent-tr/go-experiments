package store

import (
	"bytes"
	"context"
	"fmt"
	"mylife-tools-server/log"
	"mylife-tools-server/services/database"
	"mylife-tools-server/services/io"
	"mylife-tools-server/services/tasks"
	"mylife-tools-server/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type collectionUpdater[TEntity Entity] struct {
	col                *collection[TEntity]
	colChangedCallback func(event *Event[TEntity])
	dataCol            *mongo.Collection
	decode             func(raw []byte) (TEntity, error)
	encode             func(obj TEntity) ([]byte, error)
	databaseUpdating   bool
	dbContext          context.Context
	dbTerminate        context.CancelFunc
	worker             *utils.Worker
}

func makeUpdater[TEntity Entity](col *collection[TEntity], dataCol *mongo.Collection, decode func(raw []byte) (TEntity, error), encode func(obj TEntity) ([]byte, error)) *collectionUpdater[TEntity] {
	updater := &collectionUpdater[TEntity]{
		col:              col,
		dataCol:          dataCol,
		decode:           decode,
		encode:           encode,
		databaseUpdating: false,
	}

	updater.colChangedCallback = func(event *Event[TEntity]) {
		updater.colChanged(event)
	}

	col.AddListener(&updater.colChangedCallback)

	updater.dbContext, updater.dbTerminate = context.WithCancel(context.Background())
	updater.worker = utils.NewWorker(updater.dbWatcher)

	return updater
}

func (updater *collectionUpdater[TEntity]) terminate() {
	updater.col.RemoveListener(&updater.colChangedCallback)
	updater.dbTerminate()
	updater.worker.Terminate()
}

func (updater *collectionUpdater[TEntity]) colChanged(event *Event[TEntity]) {
	if updater.databaseUpdating {
		return // do not persist if we are triggered from database update
	}

	tasks.Submit(storeUpdateQueueId, fmt.Sprintf("database-updater/%s-%d", updater.col.Name(), event.Type()), func() {
		_, err := database.WithTransaction(func() (interface{}, error) {
			switch event.Type() {
			case Create:
				record, err := updater.encode(event.After())
				if err != nil {
					return nil, err
				}

				_, err = updater.dataCol.InsertOne(context.TODO(), record)
				if err != nil {
					return nil, err
				}

			case Update:
				record, err := updater.encode(event.After())
				if err != nil {
					return nil, err
				}

				filter, err := buildIdFilter(event.Before())
				if err != nil {
					return nil, err
				}

				_, err = updater.dataCol.ReplaceOne(context.TODO(), filter, record)
				if err != nil {
					return nil, err
				}

			case Remove:
				filter, err := buildIdFilter(event.Before())
				if err != nil {
					return nil, err
				}

				_, err = updater.dataCol.DeleteOne(context.TODO(), filter)

			default:
				return nil, fmt.Errorf("Unsupported event type: '%d'", event.Type())
			}

			return nil, nil
		})

		if err != nil {
			logger.WithError(err).WithFields(log.Fields{"collectionName": updater.col.Name(), "eventType": event.Type()}).Error("Could not process event")
		}
	})
}

func buildIdFilter[TEntity Entity](obj TEntity) (interface{}, error) {
	id, err := database.EncodeId(obj.Id())
	if err != nil {
		return nil, err
	}

	return bson.D{{Key: "_id", Value: id}}, nil
}

type documentKey struct {
	Id primitive.ObjectID `bson:"_id"`
}

type changeEvent struct {
	OperationType string      `bson:"operationType"`
	DocumentKey   documentKey `bson:"documentKey"`
	FullDocument  bson.Raw    `bson:"fullDocument"`
	SessionID     bson.Raw    `bson:"lsid"`
}

func (updater *collectionUpdater[TEntity]) dbWatcher(exit chan struct{}) {
	// First load existings
	initialData, err := updater.dataCol.Find(updater.dbContext, bson.D{})
	if err != nil {
		logger.WithError(err).WithField("collectionName", updater.col.Name()).Error("Could not fetch initial data")
		return
	}

	initialObjs := make([]TEntity, 0)

	for initialData.Next(updater.dbContext) {
		obj, err := updater.decode(initialData.Current)
		if err != nil {
			logger.WithError(err).WithField("collectionName", updater.col.Name()).Error("Could not decode initial data")
			continue
		}

		initialObjs = append(initialObjs, obj)
	}

	logger.WithField("collectionName", updater.col.Name()).Debugf("Fetched %d records", len(initialObjs))

	io.SubmitIoTask(fmt.Sprintf("store-updater/%s-load", updater.col.Name()), func() {
		updater.databaseUpdating = true
		updater.col.container.ReplaceAll(initialObjs, nil)
		updater.databaseUpdating = false
	})

	// Then watch (not that load is not atomic so we may miss first events?)
	changeStream, err := updater.dataCol.Watch(updater.dbContext, mongo.Pipeline{}, options.ChangeStream().SetFullDocument(options.UpdateLookup))
	if err != nil {
		logger.WithError(err).WithField("collectionName", updater.col.Name()).Error("Could not watch change stream")
		return
	}

	for changeStream.Next(updater.dbContext) {
		change := changeEvent{}
		err := bson.Unmarshal(changeStream.Current, &change)
		if err != nil {
			logger.WithError(err).WithField("collectionName", updater.col.Name()).Error("Could not unmarshal change stream event")
			continue
		}

		updater.handleDbChange(&change)
	}
}

func (updater *collectionUpdater[TEntity]) handleDbChange(change *changeEvent) {
	if change.SessionID != nil && bytes.Equal(change.SessionID, database.SessionId()) {
		logger.WithField("collectionName", updater.col.Name()).Debug("Got change with same session id than current session, ignored")
		return
	}

	switch change.OperationType {
	case "insert", "replace", "update":
		obj, err := updater.decode(change.FullDocument)
		if err != nil {
			logger.WithError(err).WithField("collectionName", updater.col.Name()).Error("Could not unmarshal object")
			return
		}

		io.SubmitIoTask(fmt.Sprintf("store-updater/%s-set", updater.col.Name()), func() {
			updater.databaseUpdating = true
			updater.col.Set(obj)
			updater.databaseUpdating = false
		})

	case "delete":
		id := database.DecodeId(change.DocumentKey.Id)

		io.SubmitIoTask(fmt.Sprintf("store-updater/%s-set", updater.col.Name()), func() {
			updater.databaseUpdating = true
			updater.col.Delete(id)
			updater.databaseUpdating = false
		})

	default:
		logger.WithField("operationType", change.OperationType).Warn("Unhandled database change stream operation type")
	}
}
