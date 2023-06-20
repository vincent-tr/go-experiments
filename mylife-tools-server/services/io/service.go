package io

import (
	"mylife-tools-server/log"
	"mylife-tools-server/services"
	"mylife-tools-server/services/sessions"
	"mylife-tools-server/services/tasks"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
)

var logger = log.CreateLogger("mylife:server:io")

func init() {
	services.Register(&ioService{})
}

type ioService struct {
	server *socketio.Server
}

func (service *ioService) Init() error {
	if err := tasks.CreateQueue("io"); err != nil {
		return err
	}

	service.server = socketio.NewServer(nil)

	service.server.OnConnect("/", func(socket socketio.Conn) error {
		session := sessions.NewSession()
		ioSession := newIoSession(session, socket)
		session.RegisterStateObject("io", ioSession)
		socket.SetContext(ioSession)
		return nil
	})

	service.server.OnEvent("/", "message", func(socket socketio.Conn, msg string) {
		ioSession := getIoSession(socket)
		ioSession.dispatch(msg)
	})

	service.server.OnError("/", func(socket socketio.Conn, err error) {
		ioSession := getIoSession(socket)
		logger.WithError(err).WithField("sessionId", ioSession.session.Id()).Error("Got error on socket")
	})

	service.server.OnDisconnect("/", func(socket socketio.Conn, reason string) {
		ioSession := getIoSession(socket)
		sessions.CloseSession(ioSession.session)
	})

	go func() {
		if err := service.server.Serve(); err != nil {
			logger.WithError(err).Error("socketio listen error")
		}
	}()

	return nil
}

func (service *ioService) Terminate() error {
	if err := service.server.Close(); err != nil {
		return err
	}

	return tasks.CloseQueue("io")
}

func (service *ioService) ServiceName() string {
	return "io"
}

func (service *ioService) Dependencies() []string {
	return []string{"api", "sessions", "tasks"}
}

func getIoSession(socket socketio.Conn) *ioSession {
	return socket.Context().(*ioSession)
}

func getService() *ioService {
	return services.GetService[*ioService]("io")
}

// Public access

func GetHandler() http.Handler {
	return getService().server
}

func NotifySession(session sessions.Session, notification any) {
	ios := session.GetStateObject("io").(*ioSession)
	ios.notify(notification)
}

func SubmitIoTask(name string, impl tasks.Task) error {
	return tasks.Submit("io", name, impl)
}
