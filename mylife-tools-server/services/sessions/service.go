package sessions

import (
	"mylife-tools-server/log"
	"mylife-tools-server/services"
	"mylife-tools-server/utils"
)

var logger = log.CreateLogger("mylife:server:sessions")

func init() {
	services.Register(&sessionService{})
}

type sessionService struct {
	sessions map[uint64]*Session
	idGen    utils.IdGenerator
}

func (service *sessionService) Init(arg interface{}) error {
	service.sessions = make(map[uint64]*Session)
	service.idGen = utils.NewIdGenerator()

	return nil
}

func (service *sessionService) Terminate() error {
	return nil
}

func (service *sessionService) ServiceName() string {
	return "sessions"
}

func (service *sessionService) Dependencies() []string {
	return []string{}
}

func (service *sessionService) newSession() *Session {
	var id = service.idGen.Next()

	var session = &Session{
		id:          id,
		onTerminate: make([]TerminateCallback, 0),
		state:       make(map[string]interface{}),
	}

	service.sessions[id] = session
	logger.WithField("sessionId", session.id).Debug("New session")
	return session
}

func (service *sessionService) closeSession(session *Session) {
	delete(service.sessions, session.id)
	session.terminate()
	logger.WithField("sessionId", session.id).Debug("Session closed")
}

func getService() *sessionService {
	return services.GetService[*sessionService]("sessions")
}

// Public access

func NewSession() *Session {
	return getService().newSession()
}

func CloseSession(session *Session) {
	getService().closeSession(session)
}
