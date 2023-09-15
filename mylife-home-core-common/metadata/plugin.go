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

type Plugin struct {
	module      string
	name        string
	description string
	usage       PluginUsage
	config      map[string]*ConfigItem
	members     map[string]*Member
}

func (plugin *Plugin) Module() string {
	return plugin.module
}

func (plugin *Plugin) Name() string {
	return plugin.name
}

func (plugin *Plugin) Id() string {
	return plugin.module + "." + plugin.name
}

func (plugin *Plugin) Description() string {
	return plugin.description
}

func (plugin *Plugin) Usage() PluginUsage {
	return plugin.usage
}

func (plugin *Plugin) ConfigNames() []string {
	return maps.Keys(plugin.config)
}

func (plugin *Plugin) Config(name string) *ConfigItem {
	return plugin.config[name]
}

func (plugin *Plugin) MemberNames() []string {
	return maps.Keys(plugin.members)
}

func (plugin *Plugin) Member(name string) *Member {
	return plugin.members[name]
}
