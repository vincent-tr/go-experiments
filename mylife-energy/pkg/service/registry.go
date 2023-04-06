package service

import (
	log "mylife-energy/pkg/log"
	"os"
	"os/signal"
	"syscall"
)

type orderData struct {
	order []string
	set   map[string]struct{}
}

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

func Init(services []string) {
	logger.WithField("services", services).Debug("Service registry init")

	data := orderData{order: []string{}, set: make(map[string]struct{})}
	for _, service := range services {
		computeInitOrder(service, &data, 0)
	}

	for _, name := range data.order {
		service := mustService(name)
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

func RunServices(services []string) {
	Init(services)
	waitSignal()
	Terminate()
}

func GetService[T any](name string) T {
	value, ok := mustService(name).(T)
	if !ok {
		logger.WithField("name", name).Fatal("Service bad type")
	}

	return value
}

func mustService(name string) Service {
	service, ok := registry[name]
	if !ok {
		logger.WithField("name", name).Fatal("Service does not exist")
	}

	return service
}

func computeInitOrder(name string, data *orderData, recursiveCount int) {
	if recursiveCount > 50 {
		logger.Fatal("Cyclic service dependency")
	}

	if _, exists := data.set[name]; exists {
		return
	}

	dependencies := mustService(name).Dependencies()
	for _, dependency := range dependencies {
		computeInitOrder(dependency, data, recursiveCount+1)
	}

	data.order = append(data.order, name)
	data.set[name] = struct{}{}
}

func waitSignal() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigs

	logger.WithField("signal", sig).Info("Got signal")
}
