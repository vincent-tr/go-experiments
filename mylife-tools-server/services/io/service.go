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
	services.Register(&ioService{})
}

type ioService struct {
}

func (service *ioService) Init() error {
	return nil
}

func (service *ioService) Terminate() error {
	return nil
}

func (service *ioService) ServiceName() string {
	return "io"
}

func (service *ioService) Dependencies() []string {
	return []string{"api", "sessions"}
}

func (service *ioService) Handler(writer http.ResponseWriter, reader *http.Request) {
	socket, err := websocket.Accept(writer, reader, nil)
	if err != nil {
		logger.WithField("error", err).Error("Accept error")
		return
	}

	session := sessions.NewSession()

	ioSession := makeSession(session, socket)
	session.RegisterStateObject("io", ioSession)
}

func getService() *ioService {
	return services.GetService[*ioService]("io")
}

// Public access

func GetHandler(name string) func(writer http.ResponseWriter, reader *http.Request) {
	return getService().Handler
}

func NotifySession(session sessions.Session, notification any) {
	ios := session.GetStateObject("io").(*ioSession)
	ios.notify(notification)
}
