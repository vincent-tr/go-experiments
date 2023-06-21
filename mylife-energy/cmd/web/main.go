package main

import (
	"mylife-tools-server/log"
	"mylife-tools-server/services"
	_ "mylife-tools-server/services/web"
)

var logger = log.CreateLogger("mylife:energy:test")

func main() {
	args := make(map[string]interface{})
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
