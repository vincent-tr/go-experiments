package notification

import (
	"mylife-tools-server/log"
	"mylife-tools-server/services"
	"mylife-tools-server/services/sessions"
	"mylife-tools-server/services/store"
)

var logger = log.CreateLogger("mylife:server:notification")

func init() {
	services.Register(&notificationService{})
}

type notificationService struct {
}

func (service *notificationService) Init(arg interface{}) error {
	return nil
}

func (service *notificationService) Terminate() error {
	return nil
}

func (service *notificationService) ServiceName() string {
	return "notification"
}

func (service *notificationService) Dependencies() []string {
	return []string{"sessions", "io"}
}

func (service *notificationService) getNotifications(session *sessions.Session) *notificationSession {
	if existing := session.FindStateObject("notification"); existing != nil {
		return existing.(*notificationSession)
	} else {
		notificationSession := newNotificationSession(session)
		session.RegisterStateObject("notification", notificationSession)
		return notificationSession
	}
}

func getService() *notificationService {
	return services.GetService[*notificationService]("notification")
}

// Public access

func NotifyView[TEntity store.Entity](session *sessions.Session, view store.IContainer[TEntity]) uint64 {
	notificationSession := getService().getNotifications(session)
	return registerView[TEntity](notificationSession, view)
}

func UnnotifyView(session *sessions.Session, viewId uint64) {
	notificationSession := getService().getNotifications(session)
	notificationSession.closeView(viewId)
}
