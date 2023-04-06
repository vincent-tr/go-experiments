package main

import (
	"context"
	"encoding/json"
	"time"

	"mylife-energy/pkg/config"
	"mylife-energy/pkg/database"
	"mylife-energy/pkg/log"
	"mylife-energy/pkg/mqtt"
	"mylife-energy/pkg/services"
	"mylife-energy/pkg/utils"
)

var logger = log.CreateLogger("main")

type MongoConfig = string

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

type MainService struct {
	records chan Record
	worker  utils.Worker
}

func (service *MainService) Init() error {
	service.records = make(chan Record, 100)
	service.worker = utils.InitWorker(func(exit chan struct{}) {
		service.workerEntry(exit)
	})

	mqtt.Subscribe("+/energy", func(topic string, data []byte) {
		service.handleMessage(topic, data)
	})

	return nil
}

func (service *MainService) Terminate() error {

	return nil
}

func (service *MainService) ServiceName() string {
	return "main"
}

func (service *MainService) Dependencies() []string {
	return []string{"mqtt", "database"}
}

func init() {
	services.Register(&MainService{})
}

func (service *MainService) handleMessage(topic string, data []byte) {
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

func (service *MainService) workerEntry(exit chan struct{}) {
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

func main() {

	mongoConfig := config.GetString("mongo")
	logger.WithField("mongoConfig", mongoConfig).Info("Config")

	services.RunServices([]string{"main"})
}
