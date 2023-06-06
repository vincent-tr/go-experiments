package sessions

import (
	"mylife-tools-server/services"
)

func init() {
	services.Register(&SessionService{})
}

type SessionService struct {
	sessions map[int]*Session
	idGen    idGenerator
}

func (service *SessionService) Init() error {
	service.sessions = make(map[int]*Session)
	service.idGen = newIdGenerator()

	return nil
}

func (service *SessionService) Terminate() error {
	return nil
}

func (service *SessionService) ServiceName() string {
	return "sessions"
}

func (service *SessionService) Dependencies() []string {
	return []string{}
}

func (service *SessionService) NewSession() *Session {
	var id = service.idGen.Next()
	var session = &Session{id: id}
	service.sessions[id] = session
	logger.WithField("sessionId", session.id).Debug("New session")
	return session
}

func (service *SessionService) CloseSession(session *Session) {
	delete(service.sessions, session.id)
	session.terminate()
	logger.WithField("sessionId", session.id).Debug("Session closed")
}
