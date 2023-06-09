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

type Message struct {
	Id                string  `json:"id"`
	DeviceClass       string  `json:"device_class"`
	StateClass        string  `json:"state_class"`
	UnitOfMeasurement string  `json:"unit_of_measurement"`
	AccuracyDecimals  int     `json:"accuracy_decimals"`
	Value             float64 `json:"value"`
}

type SensorData struct {
	SensorId          string `bson:"sensorId"`
	DeviceClass       string `bson:"deviceClass"`
	StateClass        string `bson:"stateClass"`
	UnitOfMeasurement string `bson:"unitOfMeasurement"`
	AccuracyDecimals  int    `bson:"accuracyDecimals"`
}

type Record struct {
	Timestamp time.Time  `bson:"timestamp"`
	Sensor    SensorData `bson:"sensor"`
	Value     float64    `bson:"value"`
}

type CollectorService struct {
	records chan Record
	worker  utils.Worker
}

func (service *CollectorService) Init() error {
	service.records = make(chan Record, 100)
	service.worker = utils.InitWorker(service.workerEntry)

	mqtt.Subscribe("+/energy", func(topic string, data []byte) {
		service.handleMessage(topic, data)
	})

	return nil
}

func (service *CollectorService) Terminate() error {

	return nil
}

func (service *CollectorService) ServiceName() string {
	return "collector"
}

func (service *CollectorService) Dependencies() []string {
	return []string{"mqtt", "database"}
}

func init() {
	services.Register(&CollectorService{})
}

func (service *CollectorService) handleMessage(topic string, data []byte) {
	logger.WithFields(log.Fields{"data": string(data), "topic": topic}).Debug("Got message")

	message := Message{}
	if err := json.Unmarshal(data, &message); err != nil {
		logger.WithFields(log.Fields{"error": err, "data": data}).Error("Error reading JSON")
		return
	}

	record := Record{
		Timestamp: time.Now(),
		Sensor: SensorData{
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

func (service *CollectorService) workerEntry(exit chan struct{}) {
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
