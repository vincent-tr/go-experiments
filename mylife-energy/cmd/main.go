package main

import (
	"encoding/json"
	mqtt "mylife-energy/pkg/mqtt"
	serviceRegistry "mylife-energy/pkg/service"

	config "mylife-energy/pkg/config"
	log "mylife-energy/pkg/log"
)

var logger = log.CreateLogger("main")

type MongoConfig = string

type Measure struct {
	Id                string  `json:"id"`
	DeviceClass       string  `json:"device_class"`
	StateClass        string  `json:"state_class"`
	UnitOfMeasurement string  `json:"unit_of_measurement"`
	AccuracyDecimals  int     `json:"accuracy_decimals"`
	Value             float64 `json:"value"`
}

type MainService struct {
}

func (service *MainService) Init() error {
	mqtt := serviceRegistry.GetService[*mqtt.MqttService]("mqtt")
	mqtt.Subscribe("+/energy", func(data []byte) {

		measure := Measure{}
		if err := json.Unmarshal(data, &measure); err != nil {
			logger.WithFields(log.Fields{"error": err, "data": data}).Error("Error reading JSON")
			return
		}

		logger.WithField("measure", measure).Info("Measure")
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
	return []string{"mqtt"}
}

func init() {
	serviceRegistry.Register(&MainService{})
}

func main() {

	mongoConfig := config.GetString("mongo")
	logger.WithField("mongoConfig", mongoConfig).Info("Config")

	serviceRegistry.RunServices([]string{"main"})
}
