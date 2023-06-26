package live

import (
	"mylife-tools-server/log"
	"mylife-tools-server/services"
	"mylife-tools-server/services/io"
	"mylife-tools-server/services/store"
	"mylife-tools-server/utils"
	"time"
)

var logger = log.CreateLogger("mylife:energy:live")

type liveService struct {
	worker   *utils.Worker
	measures *store.Container[*Measure]
	sensors  *store.Container[*Sensor]
}

func (service *liveService) Init(arg interface{}) error {
	service.measures = store.NewContainer[*Measure]("measures")
	service.sensors = store.NewContainer[*Sensor]("sensors")

	service.worker = utils.NewWorker(service.workerEntry)

	return nil
}

func (service *liveService) Terminate() error {
	service.worker.Terminate()

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
	results, err := fetchResults()

	if err != nil {
		logger.WithError(err).Error("error")
		return
	}

	io.SubmitIoTask("live/sync", func() {
		for _, result := range results {
			logger.Infof("Item: %+v\n", result)
		}

		service.measures.Set(&Measure{
			id:        "id1",
			timestamp: time.Now(),
			value:     42,
			sensor:    "the sensor",
		})

		service.sensors.Set(&Sensor{
			id:                "test",
			sensorId:          "sensorId",
			deviceClass:       "deviceClass",
			stateClass:        "stateClass",
			unitOfMeasurement: "unitOfMeasurement",
			accuracyDecimals:  42,
		})

		// reflect.DeepEqual

	})
}

func getService() *liveService {
	return services.GetService[*liveService]("live")
}

// Public access

func GetSensors() store.IContainer[*Sensor] {
	return getService().sensors
}

func GetMeasures() store.IContainer[*Measure] {
	return getService().measures
}
