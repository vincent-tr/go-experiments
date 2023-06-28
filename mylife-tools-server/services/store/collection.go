package store

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ICollection[TEntity Entity] interface {
	IContainer[TEntity]

	Set(obj TEntity) (TEntity, error)
	Delete(id string) bool
}

type genericInternalCollection interface {
	terminate()
}

type collection[TEntity Entity] struct {
	container *Container[TEntity]
	updater   *collectionUpdater[TEntity]
}

func (col *collection[TEntity]) AddListener(callback *func(event *Event[TEntity])) {
	col.container.AddListener(callback)
}

func (col *collection[TEntity]) RemoveListener(callback *func(event *Event[TEntity])) {
	col.container.RemoveListener(callback)
}

func (col *collection[TEntity]) Name() string {
	return col.container.Name()
}

func (col *collection[TEntity]) Find(id string) (TEntity, bool) {
	return col.container.Find(id)
}

func (col *collection[TEntity]) Get(id string) (TEntity, error) {
	return col.container.Get(id)
}

func (col *collection[TEntity]) List() []TEntity {
	return col.container.List()
}

func (col *collection[TEntity]) Size() int {
	return col.container.Size()
}

func (col *collection[TEntity]) Filter(predicate func(obj TEntity) bool) []TEntity {
	return col.container.Filter(predicate)
}

func (col *collection[TEntity]) Exists(predicate func(obj TEntity) bool) bool {
	return col.container.Exists(predicate)
}

func (col *collection[TEntity]) Set(obj TEntity) (TEntity, error) {
	oid, err := primitive.ObjectIDFromHex(obj.Id())

	if err != nil {
		var zero TEntity
		return zero, fmt.Errorf("Invalid entity id in collection '%s': %w", col.Name(), err)
	} else if oid.IsZero() {
		var zero TEntity
		return zero, fmt.Errorf("Invalid entity id in collection '%s': empty id", col.Name())
	}

	return col.container.Set(obj), nil
}

func (col *collection[TEntity]) Delete(id string) bool {
	return col.container.Delete(id)
}

func (col *collection[TEntity]) terminate() {
	col.updater.terminate()
}
