package main

import (
	"mylife-tools-server/log"
	"mylife-tools-server/services"
	"mylife-tools-server/sessions"
)

var logger = log.CreateLogger("mylife:energy:test")

func main() {
	services.RunServices([]string{"test"})
}

type TestService struct {
	session *sessions.Session
}

func (service *TestService) Init() error {
	var sessionsService = services.GetService[*sessions.SessionService]("sessions")
	service.session = sessionsService.NewSession()

	return nil
}

func (service *TestService) Terminate() error {
	var sessionsService = services.GetService[*sessions.SessionService]("sessions")
	sessionsService.CloseSession(service.session)
	service.session = nil

	return nil
}

func (service *TestService) ServiceName() string {
	return "test"
}

func (service *TestService) Dependencies() []string {
	return []string{"sessions"}
}

func init() {
	services.Register(&TestService{})
}
