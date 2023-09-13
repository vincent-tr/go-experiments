package metadata

import (
	"golang.org/x/exp/maps"
)

type PluginUsage string

const (
	Sensor   PluginUsage = "sensor"
	Actuator PluginUsage = "actuator"
	Logic    PluginUsage = "logic"
	Ui       PluginUsage = "ui"
)

type Component struct {
	name        string
	description string
	usage       PluginUsage
	config      map[string]*ConfigItem
	members     map[string]*Member
}

func (component *Component) Name() string {
	return component.name
}

func (component *Component) Description() string {
	return component.description
}

func (component *Component) Usage() PluginUsage {
	return component.usage
}

func (component *Component) ConfigNames() []string {
	return maps.Keys(component.config)
}

func (component *Component) Config(name string) *ConfigItem {
	return component.config[name]
}

func (component *Component) MemberNames() []string {
	return maps.Keys(component.members)
}

func (component *Component) Member(name string) *Member {
	return component.members[name]
}
