package main

import (
	mqtt "mylife-energy/pkg/mqtt"
	serviceRegistry "mylife-energy/pkg/service"

	config "mylife-energy/pkg/config"
	log "mylife-energy/pkg/log"
)

var logger = log.CreateLogger("main")

type MongoConfig = string

type MainService struct {
}

func (service *MainService) Init() error {
	mqtt := serviceRegistry.GetService[*mqtt.MqttService]("mqtt")
	mqtt.Subscribe("+/energy", func(data []byte) {
		logger.WithField("msg", string(data)).Info("Message")
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
