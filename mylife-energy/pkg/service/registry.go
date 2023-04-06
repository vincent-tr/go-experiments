package service

import (
	log "mylife-energy/pkg/log"
)

var registry = make(map[string]Service)
var running = []Service{}

var logger = log.CreateLogger("service:registry")

func Register(service Service) {
	name := service.ServiceName()
	if _, ok := registry[name]; ok {
		logger.WithField("name", name).Fatal("Service already registered")
	}

	registry[name] = service
	logger.WithField("name", name).Info("Service registered")
}

func Init() {
	logger.Debug("Service registry init")
	// TODO: compute deps + only init needed services

	for _, service := range registry {
		if err := service.Init(); err != nil {
			logger.WithFields(log.Fields{"name": service.ServiceName(), "error": err}).Fatal("Service init failed")
		}

		running = append(running, service)
	}
}

func Terminate() {
	logger.Debug("Service registry terminate")

	for _, service := range running {
		if err := service.Terminate(); err != nil {
			logger.WithFields(log.Fields{"name": service.ServiceName(), "error": err}).Fatal("Service terminate failed")
		}
	}
}

func GetService[T any](name string) T {
	service, ok := registry[name]
	if !ok {
		logger.WithField("name", name).Fatal("Service does not exist")
	}

	value, ok := service.(T)
	if !ok {
		logger.WithField("name", name).Fatal("Service bad type")
	}

	return value
}
