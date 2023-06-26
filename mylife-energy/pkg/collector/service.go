package collector

import (
	"context"
	"encoding/json"
	"time"

	"mylife-tools-server/log"
	"mylife-tools-server/services"
	"mylife-tools-server/services/database"
	"mylife-tools-server/services/mqtt"
	"mylife-tools-server/utils"
)

var logger = log.CreateLogger("mylife:energy:collector")

type message struct {
	Id                string  `json:"id"`
	DeviceClass       string  `json:"device_class"`
	StateClass        string  `json:"state_class"`
	UnitOfMeasurement string  `json:"unit_of_measurement"`
	AccuracyDecimals  int     `json:"accuracy_decimals"`
	Value             float64 `json:"value"`
}

type sensorData struct {
	SensorId          string `bson:"sensorId"`
	DeviceClass       string `bson:"deviceClass"`
	StateClass        string `bson:"stateClass"`
	UnitOfMeasurement string `bson:"unitOfMeasurement"`
	AccuracyDecimals  int    `bson:"accuracyDecimals"`
}

type record struct {
	Timestamp time.Time  `bson:"timestamp"`
	Sensor    sensorData `bson:"sensor"`
	Value     float64    `bson:"value"`
}

type collectorService struct {
	records chan record
	worker  *utils.Worker
}

func (service *collectorService) Init(arg interface{}) error {
	service.records = make(chan record, 100)
	service.worker = utils.NewWorker(service.workerEntry)

	mqtt.Subscribe("+/energy", func(topic string, data []byte) {
		service.handleMessage(topic, data)
	})

	return nil
}

func (service *collectorService) Terminate() error {
	service.worker.Terminate()

	return nil
}

func (service *collectorService) ServiceName() string {
	return "collector"
}

func (service *collectorService) Dependencies() []string {
	return []string{"mqtt", "database"}
}

func init() {
	services.Register(&collectorService{})
}

func (service *collectorService) handleMessage(topic string, data []byte) {
	logger.WithFields(log.Fields{"data": string(data), "topic": topic}).Debug("Got message")

	message := message{}
	if err := json.Unmarshal(data, &message); err != nil {
		logger.WithError(err).WithField("data", data).Error("Error reading JSON")
		return
	}

	record := record{
		Timestamp: time.Now(),
		Sensor: sensorData{
			SensorId:          message.Id,
			DeviceClass:       message.DeviceClass,
			StateClass:        message.StateClass,
			UnitOfMeasurement: message.UnitOfMeasurement,
			AccuracyDecimals:  message.AccuracyDecimals,
		},
		Value: message.Value,
	}

	service.records <- record
}

func (service *collectorService) workerEntry(exit chan struct{}) {
	collection := database.GetCollection("measures")

	for {
		select {
		case <-exit:
			return
		case record := <-service.records:
			logger.WithField("record", record).Debug("Insert record")
			if _, err := collection.InsertOne(context.TODO(), record); err != nil {
				logger.WithError(err).Error("Error inserting record")
			}
		}
	}
}
