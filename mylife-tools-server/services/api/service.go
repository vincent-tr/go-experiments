package api

import (
	"fmt"
	"mylife-tools-server/log"
	"mylife-tools-server/services"
)

var logger = log.CreateLogger("mylife:server:api")

func init() {
	services.Register(&apiService{})
}

type serviceImpl struct {
	name    string
	methods map[string]*Method
}

type apiService struct {
	services map[string]*serviceImpl
}

func (service *apiService) Init(arg interface{}) error {
	service.services = make(map[string]*serviceImpl)

	if arg != nil {
		defs := arg.([]ServiceDefinition)

		for _, def := range defs {
			service.RegisterService(def)
		}
	}

	return nil
}

func (service *apiService) Terminate() error {
	for serviceName, _ := range service.services {
		delete(service.services, serviceName)
		logger.WithField("serviceName", serviceName).Info("Service unregistered")
	}

	return nil
}

func (service *apiService) Lookup(serviceName string, methodName string) (*Method, error) {
	svc, ok := service.services[serviceName]

	if !ok {
		return nil, fmt.Errorf("Service '%s' does not exist", serviceName)
	}

	method, ok := svc.methods[methodName]

	if !ok {
		return nil, fmt.Errorf("Method '%s' does not exist on service '%s'", methodName, serviceName)
	}

	return method, nil
}

func (service *apiService) RegisterService(def ServiceDefinition) {
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

func (service *apiService) ServiceName() string {
	return "api"
}

func (service *apiService) Dependencies() []string {
	return []string{}
}

func getService() *apiService {
	return services.GetService[*apiService]("api")
}

// Public access

func RegisterService(def ServiceDefinition) {
	getService().RegisterService(def)
}

func Lookup(serviceName string, methodName string) (*Method, error) {
	return getService().Lookup(serviceName, methodName)
}
