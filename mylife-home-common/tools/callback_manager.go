package tools

import "sync"

type RegistrationToken int

type CallbackRegistration[TArg any] interface {
	Register(callback func(TArg)) RegistrationToken
	Unregister(token RegistrationToken)
}

type CallbackManager[TArg any] struct {
	callbacks     map[RegistrationToken]func(TArg)
	callbacksSync sync.RWMutex
	nextToken     RegistrationToken
}

func NewCallbackManager[TArg any]() *CallbackManager[TArg] {
	return &CallbackManager[TArg]{
		callbacks: make(map[RegistrationToken]func(TArg)),
		nextToken: 1,
	}
}

func (m *CallbackManager[TArg]) Execute(arg TArg) {
	m.callbacksSync.RLock()
	defer m.callbacksSync.RUnlock()

	for _, callback := range m.callbacks {
		callback(arg)
	}
}

func (m *CallbackManager[TArg]) Register(callback func(TArg)) RegistrationToken {
	m.callbacksSync.Lock()
	defer m.callbacksSync.Unlock()

	token := m.nextToken
	m.nextToken += 1

	m.callbacks[token] = callback

	return token
}

func (m *CallbackManager[TArg]) Unregister(token RegistrationToken) {
	m.callbacksSync.Lock()
	defer m.callbacksSync.Unlock()

	delete(m.callbacks, token)
}
