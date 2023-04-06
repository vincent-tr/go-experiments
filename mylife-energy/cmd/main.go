package main

import (
	"time"

	mqtt "mylife-energy/pkg/mqtt"
	serviceRegistry "mylife-energy/pkg/service"

	config "mylife-energy/pkg/config"
	log "mylife-energy/pkg/log"
)

var logger = log.CreateLogger("main")

type MongoConfig = string

func main() {

	mongoConfig := config.GetString("mongo")

	logger.WithField("mongoConfig", mongoConfig).Info("Config")

	serviceRegistry.Init()

	mqtt := serviceRegistry.GetService[*mqtt.Mqtt]("mqtt")
	mqtt.Subscribe("+/energy", func(data []byte) {
		logger.WithField("msg", string(data)).Info("Message")
	})

	time.Sleep(30 * time.Second)

	serviceRegistry.Terminate()
}
