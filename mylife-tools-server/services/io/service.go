package io

import (
	"context"
	"mylife-tools-server/log"
	"mylife-tools-server/services"
	"mylife-tools-server/services/sessions"
	"mylife-tools-server/services/tasks"
	"net/http"

	socketio "github.com/vchitai/go-socket.io/v4"
)

var logger = log.CreateLogger("mylife:server:io")

func init() {
	services.Register(&ioService{})
}

type ioService struct {
	server *socketio.Server
}

func (service *ioService) Init(arg interface{}) error {
	if err := tasks.CreateQueue("io"); err != nil {
		return err
	}

	service.server = socketio.NewServer(nil)

	service.server.OnConnect("/", func(socket socketio.Conn, dispatchData map[string]interface{}) error {
		session := sessions.NewSession()
		ioSession := newIoSession(session, socket)
		session.RegisterStateObject("io", ioSession)
		socket.SetContext(context.WithValue(context.TODO(), "ioSession", ioSession))

		logger.WithFields(log.Fields{"sessionId": ioSession.session.Id(), "socketId": socket.ID()}).Debug("New IO session")

		return nil
	})

	service.server.OnEvent("/", "message", func(socket socketio.Conn, data map[string]interface{}) {
		ioSession := getIoSession(socket)

		if ioSession == nil {
			logger.WithField("socketId", socket.ID()).Error("Got message but could not get associoated session")
			return
		}

		ioSession.dispatch(data)
	})

	service.server.OnError("/", func(socket socketio.Conn, err error) {
		ioSession := getIoSession(socket)

		fields := log.Fields{}
		if socket != nil {
			fields["socketId"] = socket.ID()
		}

		if ioSession != nil {
			fields["sessionId"] = ioSession.session.Id()
		}

		logger.WithError(err).WithFields(fields).Error("Got error on socket")
	})

	service.server.OnDisconnect("/", func(socket socketio.Conn, reason string, dispatchData map[string]interface{}) {
		ioSession := getIoSession(socket)

		if ioSession == nil {
			logger.WithField("socketId", socket.ID()).Error("Got socket disconnection but could not get associoated session")
			return
		}

		logger.WithField("sessionId", ioSession.session.Id()).Debug("IO session disconnected")

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
	if socket == nil {
		return nil
	}

	ios, ok := socket.Context().Value("ioSession").(*ioSession)
	if ok {
		return ios
	} else {
		return nil
	}
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
