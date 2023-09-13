package metadata

import "github.com/gookit/goutil/errorx/panics"

type ComponentBuilder struct {
	target *Component
}

func MakeBuilder(name string, description string, usage PluginUsage) *ComponentBuilder {
	return &ComponentBuilder{
		target: &Component{
			name:        name,
			description: description,
			usage:       usage,
			config:      make(map[string]*ConfigItem),
			members:     make(map[string]*Member),
		},
	}
}

func (builder *ComponentBuilder) AddConfig(name string, description string, valueType ConfigType) *ComponentBuilder {
	_, exists := builder.target.config[name]
	panics.IsFalse(exists)

	builder.target.config[name] = &ConfigItem{name, description, valueType}

	return builder
}

func (builder *ComponentBuilder) AddState(name string, description string, valueType Type) *ComponentBuilder {
	_, exists := builder.target.members[name]
	panics.IsFalse(exists)

	builder.target.members[name] = &Member{name, description, State, valueType}

	return builder
}

func (builder *ComponentBuilder) AddAction(name string, description string, valueType Type) *ComponentBuilder {
	_, exists := builder.target.members[name]
	panics.IsFalse(exists)

	builder.target.members[name] = &Member{name, description, Action, valueType}

	return builder
}
