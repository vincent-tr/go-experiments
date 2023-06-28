package store

import (
	"fmt"

	"golang.org/x/exp/maps"
)

type EventType int

const (
	Create EventType = iota
	Update
	Remove
)

type Event[TEntity Entity] struct {
	typ    EventType
	before TEntity
	after  TEntity
}

func (event *Event[TEntity]) Type() EventType {
	return event.typ
}

func (event *Event[TEntity]) Before() TEntity {
	return event.before
}

func (event *Event[TEntity]) After() TEntity {
	return event.after
}

type IContainer[TEntity Entity] interface {
	IEventEmitter[Event[TEntity]]

	Name() string
	Find(id string) (TEntity, bool)
	Get(id string) (TEntity, error)
	List() []TEntity
	Size() int
	Filter(predicate func(obj TEntity) bool) []TEntity
	Exists(predicate func(obj TEntity) bool) bool
}

type Container[TEntity Entity] struct {
	EventEmitter[Event[TEntity]]

	name  string
	items map[string]TEntity
}

func NewContainer[TEntity Entity](name string) *Container[TEntity] {
	return &Container[TEntity]{
		EventEmitter: *NewEventEmitter[Event[TEntity]](),
		name:         name,
		items:        make(map[string]TEntity),
	}
}

func (container *Container[TEntity]) Name() string {
	return container.name
}

func (container *Container[TEntity]) Reset() {
	for k := range container.items {
		delete(container.items, k)
	}
}

func (container *Container[TEntity]) Set(obj TEntity) TEntity {
	id := obj.Id()

	existing, exists := container.items[id]
	if exists && &existing == &obj {
		// if same, no replacement, no emitted event
		return existing
	}

	container.items[id] = obj

	event := &Event[TEntity]{typ: Create, after: obj}
	if exists {
		event.typ = Update
		event.before = existing
	}

	container.Emit(event)

	return obj
}

func (container *Container[TEntity]) Delete(id string) bool {
	existing, exists := container.items[id]
	if !exists {
		return false
	}

	delete(container.items, id)

	container.Emit(&Event[TEntity]{typ: Remove, before: existing})

	return true
}

// provide nil equals to not check for equality
func (container *Container[TEntity]) ReplaceAll(objs []TEntity, equals func(a TEntity, b TEntity) bool) {
	removeSet := make(map[string]struct{})

	for id := range container.items {
		removeSet[id] = struct{}{}
	}

	for _, obj := range objs {
		delete(removeSet, obj.Id())
	}

	for id := range removeSet {
		container.Delete(id)
	}

	for _, obj := range objs {
		if container.needSet(obj, equals) {
			container.Set(obj)
		}
	}
}

func (container *Container[TEntity]) needSet(obj TEntity, equals func(a TEntity, b TEntity) bool) bool {
	if equals == nil {
		return true
	}

	existing, exists := container.Find(obj.Id())
	if !exists {
		return true
	}

	return !equals(existing, obj)
}

func (container *Container[TEntity]) Find(id string) (TEntity, bool) {
	obj, exists := container.items[id]
	return obj, exists
}

func (container *Container[TEntity]) Get(id string) (TEntity, error) {
	obj, exists := container.items[id]
	if exists {
		return obj, nil
	} else {
		return obj, fmt.Errorf("Object with id '%s' not found on collection '%s'", id, container.name)
	}
}

func (container *Container[TEntity]) List() []TEntity {
	return maps.Values(container.items)
}

func (container *Container[TEntity]) Size() int {
	return len(container.items)
}

func (container *Container[TEntity]) Filter(predicate func(obj TEntity) bool) []TEntity {
	result := make([]TEntity, 0)

	for _, value := range container.items {
		if predicate(value) {
			result = append(result, value)
		}
	}

	return result
}

func (container *Container[TEntity]) Exists(predicate func(obj TEntity) bool) bool {
	for _, value := range container.items {
		if predicate(value) {
			return true
		}
	}

	return false
}
