package live

import (
	"context"
	"mylife-energy/pkg/entities"
	"mylife-energy/pkg/services/query"
	"mylife-tools-server/log"
	"mylife-tools-server/services"
	"mylife-tools-server/services/io"
	"mylife-tools-server/services/store"
	"mylife-tools-server/utils"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

var logger = log.CreateLogger("mylife:energy:live")

type liveService struct {
	worker      *utils.Worker
	dbContext   context.Context
	dbTerminate context.CancelFunc
	measures    *store.Container[*entities.Measure]
	sensors     *store.Container[*entities.Sensor]
	pendingSync *sync.WaitGroup
}

func (service *liveService) Init(arg interface{}) error {
	service.measures = store.NewContainer[*entities.Measure]("measures")
	service.sensors = store.NewContainer[*entities.Sensor]("sensors")
	service.pendingSync = &sync.WaitGroup{}
	service.dbContext, service.dbTerminate = context.WithCancel(context.Background())

	service.worker = utils.NewWorker(service.workerEntry)

	return nil
}

func (service *liveService) Terminate() error {
	service.dbTerminate()
	service.worker.Terminate()
	service.pendingSync.Wait()

	service.measures = nil
	service.sensors = nil

	return nil
}

func (service *liveService) ServiceName() string {
	return "live"
}

func (service *liveService) Dependencies() []string {
	// io because we use io queue
	return []string{"query", "io", "tasks"}
}

func init() {
	services.Register(&liveService{})
}

func (service *liveService) workerEntry(exit chan struct{}) {

	for {
		select {
		case <-exit:
			return
		case <-time.After(10 * time.Second):
			service.sync()
		}
	}
}

func (service *liveService) sync() {
	results, err := query.Aggregate(service.dbContext, []bson.M{
		{"$sort": bson.M{"sensor.sensorId": 1, "timestamp": -1}},
		{"$group": bson.M{"_id": "$sensor.sensorId", "timestamp": bson.M{"$first": "$timestamp"}, "value": bson.M{"$first": "$value"}, "sensor": bson.M{"$first": "$sensor"}}},
	})

	if err != nil {
		logger.WithError(err).Error("Error fetching results")
		return
	}

	service.pendingSync.Add(1)

	err = io.SubmitIoTask("live/sync", func() {
		syncEntity[*entities.Measure](service.measures, results, accessMeasure, entities.MeasuresEqual)
		syncEntity[*entities.Sensor](service.sensors, results, accessSensor, entities.SensorsEqual)
		service.pendingSync.Done()
	})

	if err != nil {
		service.pendingSync.Done()

		logger.WithError(err).Error("Error submitting io task")
		return
	}
}

func accessMeasure(result *query.Result) *entities.Measure {
	return result.Measure
}

func accessSensor(result *query.Result) *entities.Sensor {
	return result.Sensor
}

func syncEntity[TEntity store.Entity](container *store.Container[TEntity], results []query.Result, access func(result *query.Result) TEntity, equals func(a TEntity, b TEntity) bool) {

	removeSet := make(map[string]struct{})

	for _, obj := range container.List() {
		removeSet[obj.Id()] = struct{}{}
	}

	for _, result := range results {
		obj := access(&result)
		delete(removeSet, obj.Id())
	}

	for id := range removeSet {
		container.Delete(id)
	}

	for _, result := range results {
		obj := access(&result)
		existing, exists := container.Find(obj.Id())

		if exists && equals(obj, existing) {
			continue
		}

		container.Set(obj)
	}
}

func getService() *liveService {
	return services.GetService[*liveService]("live")
}

// Public access

func GetSensors() store.IContainer[*entities.Sensor] {
	return getService().sensors
}

func GetMeasures() store.IContainer[*entities.Measure] {
	return getService().measures
}
