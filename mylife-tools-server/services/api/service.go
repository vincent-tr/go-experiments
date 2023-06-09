package api

import (
	"errors"
	"fmt"
	"mylife-tools-server/log"
	"mylife-tools-server/services"
)

var logger = log.CreateLogger("mylife:server:api")

func init() {
	services.Register(&ApiService{})
}

type serviceImpl struct {
	name    string
	methods map[string]*Method
}

type ApiService struct {
	services map[string]*serviceImpl
}

func (service *ApiService) Init() error {
	service.services = make(map[string]*serviceImpl)

	return nil
}

func (service *ApiService) Terminate() error {
	for serviceName, _ := range service.services {
		delete(service.services, serviceName)
		logger.WithField("serviceName", serviceName).Info("Service unregistered")
	}

	return nil
}

func (service *ApiService) Lookup(serviceName string, methodName string) (*Method, error) {
	svc, ok := service.services[serviceName]

	if !ok {
		return nil, errors.New(fmt.Sprintf("Service '%s' does not exist", serviceName))
	}

	method, ok := svc.methods[methodName]

	if !ok {
		return nil, errors.New(fmt.Sprintf("Method '%s does not exist on service '%s'", methodName, serviceName))
	}

	return method, nil
}

func (service *ApiService) RegisterService(def ServiceDefinition) {
	if _, ok := service.services[def.Name]; ok {
		logger.WithField("serviceName", def.Name).Fatal("Service already exists")
	}

	svc := &serviceImpl{name: def.Name, methods: make(map[string]*Method)}

	for methodName, callee := range def.Methods {
		svc.methods[methodName] = newMethod(callee)
	}

	service.services[svc.name] = svc

	logger.WithField("serviceName", svc.name).Info("Service registered")
}

func (service *ApiService) ServiceName() string {
	return "api"
}

func (service *ApiService) Dependencies() []string {
	return []string{}
}
