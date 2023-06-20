package io

import (
	"mylife-tools-server/log"
	"mylife-tools-server/services"
	"mylife-tools-server/services/sessions"
	"mylife-tools-server/services/tasks"
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
	return tasks.CreateQueue("io")
}

func (service *ioService) Terminate() error {
	return tasks.CloseQueue("io")
}

func (service *ioService) ServiceName() string {
	return "io"
}

func (service *ioService) Dependencies() []string {
	return []string{"api", "sessions", "tasks"}
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

func SubmitIoTask(name string, impl tasks.Task) error {
	return tasks.Submit("io", name, impl)
}
