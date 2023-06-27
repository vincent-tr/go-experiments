package live

import (
	"context"
	"mylife-energy/pkg/entities"
	"mylife-tools-server/log"
	"mylife-tools-server/services"
	"mylife-tools-server/services/io"
	"mylife-tools-server/services/store"
	"mylife-tools-server/utils"
	"sync"
	"time"
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
	return []string{"database", "io", "tasks"}
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
	results, err := fetchResults(service.dbContext)

	if err != nil {
		logger.WithError(err).Error("Error fetching results")
		return
	}

	newMeasures := make([]*entities.Measure, 0)
	newSensors := make([]*entities.Sensor, 0)

	for _, result := range results {
		newMeasure := makeMeasureFromData(&result)
		newMeasures = append(newMeasures, newMeasure)

		newSensor := makeSensorFromData(&result.Sensor)
		newSensors = append(newSensors, newSensor)
	}

	service.pendingSync.Add(1)

	err = io.SubmitIoTask("live/sync", func() {
		syncEntity[*entities.Measure](service.measures, newMeasures, entities.MeasuresEqual)
		syncEntity[*entities.Sensor](service.sensors, newSensors, entities.SensorsEqual)
		service.pendingSync.Done()
	})

	if err != nil {
		service.pendingSync.Done()

		logger.WithError(err).Error("Error submitting io task")
		return
	}
}

func syncEntity[TEntity store.Entity](container *store.Container[TEntity], list []TEntity, equals func(a TEntity, b TEntity) bool) {

	removeSet := make(map[string]struct{})

	for _, obj := range container.List() {
		removeSet[obj.Id()] = struct{}{}
	}

	for _, obj := range list {
		delete(removeSet, obj.Id())
	}

	for id := range removeSet {
		container.Delete(id)
	}

	for _, obj := range list {
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
