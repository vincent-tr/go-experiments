package store

type IEventEmitter[TEvent any] interface {
	AddListener(callback *func(event *TEvent))
	RemoveListener(callback *func(event *TEvent))
}

type EventEmitter[TEvent any] struct {
	callbacks map[*func(event *TEvent)]struct{}
}

func NewEventEmitter[TEvent any]() *EventEmitter[TEvent] {
	return &EventEmitter[TEvent]{callbacks: make(map[*func(event *TEvent)]struct{})}
}

func (emitter *EventEmitter[TEvent]) AddListener(callback *func(event *TEvent)) {
	emitter.callbacks[callback] = struct{}{}
}

func (emitter *EventEmitter[TEvent]) RemoveListener(callback *func(event *TEvent)) {
	delete(emitter.callbacks, callback)
}

func (emitter *EventEmitter[TEvent]) Emit(event *TEvent) {
	for callback := range emitter.callbacks {
		(*callback)(event)
	}
}
