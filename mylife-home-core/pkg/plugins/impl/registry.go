package impl

import (
	"reflect"
)

type RegistryImpl struct {
	types []reflect.Type
}

func (this *RegistryImpl) AddPlugin(pluginType reflect.Type) {
	this.types = append(this.types, pluginType)
}

func (this *RegistryImpl) GetPlugins() []reflect.Type {
	return this.types
}

func MakeRegistry() *RegistryImpl {
	return &RegistryImpl{
		types: make([]reflect.Type, 0),
	}
}
