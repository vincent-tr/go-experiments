package plugins

import "reflect"

type Plugin struct {
}

func BuildPlugin(pluginType reflect.Type) (*Plugin, error) {
}

func (this *Plugin) Instantiate(config map[string]any) (*Component, error) {
}

// go plugin -> module

// metadata
