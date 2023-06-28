package live

import (
	"context"
	"mylife-energy/pkg/entities"
	"mylife-energy/pkg/services/query"
	"mylife-tools-server/services/io"
	"mylife-tools-server/services/store"
	"mylife-tools-server/utils"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type fetcher struct {
	worker      *utils.Worker
	dbContext   context.Context
	dbTerminate context.CancelFunc
	measures    *store.Container[*entities.Measure]
	sensors     *store.Container[*entities.Sensor]
	pendingSync *sync.WaitGroup
}

func makeFetcher() *fetcher {
	f := &fetcher{
		measures:    store.NewContainer[*entities.Measure]("measures"),
		sensors:     store.NewContainer[*entities.Sensor]("sensors"),
		pendingSync: &sync.WaitGroup{},
	}

	f.dbContext, f.dbTerminate = context.WithCancel(context.Background())

	f.worker = utils.NewWorker(f.workerEntry)

	return f
}

func (f *fetcher) terminate() {
	f.dbTerminate()
	f.worker.Terminate()
	f.pendingSync.Wait()

	f.measures = nil
	f.sensors = nil
}

func (f *fetcher) workerEntry(exit chan struct{}) {

	for {
		select {
		case <-exit:
			return
		case <-time.After(10 * time.Second):
			f.sync()
		}
	}
}

func (f *fetcher) sync() {
	results, err := query.Aggregate(f.dbContext, []bson.M{
		{"$sort": bson.M{"sensor.sensorId": 1, "timestamp": -1}},
		{"$group": bson.M{"_id": "$sensor.sensorId", "timestamp": bson.M{"$first": "$timestamp"}, "value": bson.M{"$first": "$value"}, "sensor": bson.M{"$first": "$sensor"}}},
	})

	if err != nil {
		logger.WithError(err).Error("Error fetching results")
		return
	}

	f.pendingSync.Add(1)

	err = io.SubmitIoTask("live/fetch", func() {
		newMeasures := make([]*entities.Measure, 0, len(results))
		newSensors := make([]*entities.Sensor, 0, len(results))

		for _, result := range results {
			newMeasures = append(newMeasures, result.Measure)
			newSensors = append(newSensors, result.Sensor)
		}

		f.measures.ReplaceAll(newMeasures, entities.MeasuresEqual)
		f.sensors.ReplaceAll(newSensors, entities.SensorsEqual)

		f.pendingSync.Done()
	})

	if err != nil {
		f.pendingSync.Done()

		logger.WithError(err).Error("Error submitting io task")
		return
	}
}
