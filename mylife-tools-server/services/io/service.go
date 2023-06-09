package io

import (
	"mylife-tools-server/log"
	"mylife-tools-server/services"
	"mylife-tools-server/services/sessions"
	"net/http"

	"nhooyr.io/websocket"
)

var logger = log.CreateLogger("mylife:server:io")

func init() {
	services.Register(&IOService{})
}

type IOService struct {
}

func (service *IOService) Init() error {
	return nil
}

func (service *IOService) Terminate() error {
	return nil
}

func (service *IOService) ServiceName() string {
	return "io"
}

func (service *IOService) Dependencies() []string {
	return []string{"api", "sessions"}
}

func (service *IOService) Handler(witer http.ResponseWriter, reader *http.Request) {
	socket, err := websocket.Accept(witer, reader, nil)
	if err != nil {
		logger.WithField("error", err).Error("Accept error")
		return
	}

	sessionService := services.GetService[sessions.SessionService]("sessions")
	session := sessionService.NewSession()

	ioSession := makeSession(session, socket)
	session.RegisterStateObject("io", ioSession)
}
