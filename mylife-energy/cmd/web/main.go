package main

import (
	"mylife-tools-server/log"
	"mylife-tools-server/services"
	_ "mylife-tools-server/services/io"
	"mylife-tools-server/services/sessions"
)

var logger = log.CreateLogger("mylife:energy:test")

func main() {
	services.RunServices([]string{"test"})
}

type TestService struct {
	session *sessions.Session
}

func (service *TestService) Init() error {
	sessions.NewSession()
	service.session = sessions.NewSession()

	return nil
}

func (service *TestService) Terminate() error {
	sessions.CloseSession(service.session)
	service.session = nil

	return nil
}

func (service *TestService) ServiceName() string {
	return "test"
}

func (service *TestService) Dependencies() []string {
	return []string{"sessions", "io"}
}

func init() {
	services.Register(&TestService{})
}
