package registry

import (
	"reflect"
)

type Registry interface {
	AddPlugin(pluginType reflect.Type)
}

var registry Registry

func RegisterPlugin[T any]() {
	var ptr *T
	registry.AddPlugin(reflect.TypeOf(ptr).Elem())
}

func SetRegistry(value Registry) {
	registry = value
}
