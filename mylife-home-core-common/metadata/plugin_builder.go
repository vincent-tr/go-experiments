package metadata

import "github.com/gookit/goutil/errorx/panics"

type PluginBuilder struct {
	target *Plugin
}

func MakeBuilder(name string, description string, usage PluginUsage) *PluginBuilder {
	return &PluginBuilder{
		target: &Plugin{
			name:        name,
			description: description,
			usage:       usage,
			config:      make(map[string]*ConfigItem),
			members:     make(map[string]*Member),
		},
	}
}

func (builder *PluginBuilder) AddConfig(name string, description string, valueType ConfigType) *PluginBuilder {
	_, exists := builder.target.config[name]
	panics.IsFalse(exists)

	builder.target.config[name] = &ConfigItem{name, description, valueType}

	return builder
}

func (builder *PluginBuilder) AddState(name string, description string, valueType Type) *PluginBuilder {
	_, exists := builder.target.members[name]
	panics.IsFalse(exists)

	builder.target.members[name] = &Member{name, description, State, valueType}

	return builder
}

func (builder *PluginBuilder) AddAction(name string, description string, valueType Type) *PluginBuilder {
	_, exists := builder.target.members[name]
	panics.IsFalse(exists)

	builder.target.members[name] = &Member{name, description, Action, valueType}

	return builder
}
