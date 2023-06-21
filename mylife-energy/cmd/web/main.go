package main

import (
	"errors"
	"mylife-tools-server/log"
	"mylife-tools-server/services"
	"mylife-tools-server/services/api"
	"mylife-tools-server/services/sessions"
	_ "mylife-tools-server/services/web"
)

var logger = log.CreateLogger("mylife:energy:test")

func main() {
	args := make(map[string]interface{})

	args["api"] = []api.ServiceDefinition{
		api.MakeDefinition("common", notifySensors),
	}

	services.RunServices([]string{"test", "web"}, args)
}

type testService struct {
}

func (service *testService) Init(arg interface{}) error {
	return nil
}

func (service *testService) Terminate() error {
	return nil
}

func (service *testService) ServiceName() string {
	return "test"
}

func (service *testService) Dependencies() []string {
	return []string{}
}

func init() {
	services.Register(&testService{})
}

func notifySensors(session *sessions.Session, in struct{}) (string, error) {
	return "", errors.New("TODO")
}
