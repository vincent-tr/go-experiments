package store

import (
	"errors"
	"fmt"

	"golang.org/x/exp/maps"
)

type EventType int

const (
	Create EventType = iota
	Update
	Remove
)

type EntityConstraint interface {
	Entity
	struct{}
}

type Event[TEntity EntityConstraint] struct {
	typ    EventType
	before *TEntity
	after  *TEntity
}

func (event *Event[TEntity]) Type() EventType {
	return event.typ
}

func (event *Event[TEntity]) Before() *TEntity {
	return event.before
}

func (event *Event[TEntity]) After() *TEntity {
	return event.after
}

type IContainer[TEntity EntityConstraint] interface {
	IEventEmitter[Event[TEntity]]

	Name() string
	Find(id string) (*TEntity, bool)
	Get(id string) (*TEntity, error)
	List() []*TEntity
	Size() int
	Exists(predicate func(obj *TEntity) bool) bool
}

type Container[TEntity EntityConstraint] struct {
	name    string
	items   map[string]*TEntity
	emitter *EventEmitter[Event[TEntity]]
}

func NewContainer[TEntity EntityConstraint](name string) *Container[TEntity] {
	return &Container[TEntity]{
		name:    name,
		items:   make(map[string]*TEntity),
		emitter: NewEventEmitter[Event[TEntity]](),
	}
}

func (container *Container[TEntity]) Name() string {
	return container.name
}

// protected
func (container *Container[TEntity]) Reset() {
	for k := range container.items {
		delete(container.items, k)
	}
}

// protected
func (container *Container[TEntity]) Set(obj *TEntity) *TEntity {
	id := (*obj).Id()

	existing, exists := container.items[id]
	if exists && *existing == *obj {
		// if same, no replacement, no emitted event
		return existing
	}

	container.items[id] = obj

	event := &Event[TEntity]{typ: Create, after: obj}
	if exists {
		event.typ = Update
		event.before = existing
	}

	container.emitter.Emit(event)

	return obj
}

// protected
func (container *Container[TEntity]) Delete(id string) bool {
	existing, exists := container.items[id]
	if !exists {
		return false
	}

	delete(container.items, id)

	container.emitter.Emit(&Event[TEntity]{typ: Remove, before: existing})

	return true
}

// protected
func (container *Container[TEntity]) ReplaceAll(objs []*TEntity) {
	removeSet := make(map[string]struct{})

	for id, _ := range container.items {
		removeSet[id] = struct{}{}
	}

	for _, obj := range objs {
		delete(removeSet, (*obj).Id())
	}

	for id, _ := range removeSet {
		container.Delete(id)
	}

	for _, obj := range objs {
		container.Set(obj)
	}
}

func (container *Container[TEntity]) Find(id string) (*TEntity, bool) {
	obj, exists := container.items[id]
	return obj, exists
}

func (container *Container[TEntity]) Get(id string) (*TEntity, error) {
	obj, exists := container.items[id]
	if exists {
		return obj, nil
	} else {
		return obj, errors.New(fmt.Sprintf("Object with id '%s' not found on collection '%s'", id, container.name))
	}
}

func (container *Container[TEntity]) List() []*TEntity {
	return maps.Values(container.items)
}

func (container *Container[TEntity]) Size() int {
	return len(container.items)
}

func (container *Container[TEntity]) Filter(predicate func(obj *TEntity) bool) []*TEntity {
	result := make([]*TEntity, 0)

	for _, value := range container.items {
		if predicate(value) {
			result = append(result, value)
		}
	}

	return result
}

func (container *Container[TEntity]) Exists(predicate func(obj *TEntity) bool) bool {
	for _, value := range container.items {
		if predicate(value) {
			return true
		}
	}

	return false
}