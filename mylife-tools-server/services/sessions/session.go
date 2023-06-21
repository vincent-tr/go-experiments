package sessions

import (
	"mylife-tools-server/log"
)

type TerminateCallback func()

type SessionStateObject interface{}

type Closable interface {
	Close()
}

type Session struct {
	id          uint64
	onTerminate []TerminateCallback
	state       map[string]interface{}
}

func (session *Session) Id() uint64 {
	return session.id
}

func (session *Session) RegisterTerminateCallback(callback TerminateCallback) {
	session.onTerminate = append(session.onTerminate, callback)
}

func (session *Session) RegisterClosable(closable Closable) {
	session.RegisterTerminateCallback(closable.Close)
}

func (session *Session) RegisterStateObject(id string, state SessionStateObject) {
	session.state[id] = state

	if closable, ok := state.(Closable); ok {
		session.RegisterClosable(closable)
	}
}

func (session *Session) FindStateObject(id string) SessionStateObject {
	object, ok := session.state[id]
	if ok {
		return object
	} else {
		return nil
	}
}

func (session *Session) GetStateObject(id string) SessionStateObject {
	object, ok := session.state[id]
	if !ok {
		logger.WithFields(log.Fields{"id": id, "sessionId": session.id}).Fatal("State object does not exist")
	}

	return object
}

func (session *Session) terminate() {
	for _, callback := range session.onTerminate {
		callback()
	}

	session.onTerminate = nil
}
