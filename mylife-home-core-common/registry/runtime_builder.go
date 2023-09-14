package registry

import (
	"mylife-home-core-common/definitions"
	"mylife-home-core-common/metadata"
	"reflect"

	"github.com/gookit/goutil/errorx/panics"
)

type PluginTypeBuilder struct {
	target      *PluginType
	metaBuilder *metadata.PluginBuilder
	state       []NamedItem[*StateType]
	actions     []NamedItem[*ActionType]
	config      []NamedItem[*ConfigType]
}

type NamedItem[T any] struct {
	name   string // metadata name
	target T
}

func MakePluginTypeBuilder[TPlugin definitions.Plugin](name string, description string, usage metadata.PluginUsage) *PluginTypeBuilder {
	var ptr *TPlugin = nil
	typ := reflect.TypeOf(ptr).Elem()
	target := &PluginType{
		target: typ,
	}

	metaBuilder := metadata.MakeBuilder(name, description, usage)

	return &PluginTypeBuilder{
		target:      target,
		metaBuilder: metaBuilder,
		state:       make([]NamedItem[*StateType], 0),
		actions:     make([]NamedItem[*ActionType], 0),
		config:      make([]NamedItem[*ConfigType], 0),
	}
}

func (builder *PluginTypeBuilder) AddState(fieldName string, name string, description string, valueType metadata.Type) *PluginTypeBuilder {
	builder.metaBuilder.AddState(name, description, valueType)

	field, ok := builder.target.target.FieldByName(fieldName)
	panics.IsTrue(ok)

	stateItem := &StateType{
		target: &field,
	}

	builder.state = append(builder.state, NamedItem[*StateType]{
		name:   name,
		target: stateItem,
	})

	return builder
}

func (builder *PluginTypeBuilder) AddAction(funcImpl any, name string, description string, valueType metadata.Type) *PluginTypeBuilder {
	builder.metaBuilder.AddAction(name, description, valueType)

	funcType := reflect.TypeOf(funcImpl)
	panics.IsTrue(funcType.Kind() == reflect.Func)

	action := &ActionType{
		target: funcType,
	}

	builder.actions = append(builder.actions, NamedItem[*ActionType]{
		name:   name,
		target: action,
	})

	return builder
}

func (builder *PluginTypeBuilder) AddConfig(fieldName string, name string, description string, valueType metadata.ConfigType) *PluginTypeBuilder {
	builder.metaBuilder.AddConfig(name, description, valueType)

	field, ok := builder.target.target.FieldByName(fieldName)
	panics.IsTrue(ok)

	configItem := &ConfigType{
		target: &field,
	}

	builder.config = append(builder.config, NamedItem[*ConfigType]{
		name:   name,
		target: configItem,
	})

	return builder
}

func (builder *PluginTypeBuilder) Build() *PluginType {
	meta := builder.metaBuilder.Build()
	plugin := builder.target

	// Associate meta

	for _, stateItem := range builder.state {
		stateType := stateItem.target
		stateType.meta = meta.Member(stateItem.name)
		panics.IsTrue(stateType.meta != nil)
		panics.IsTrue(stateType.meta.MemberType() == metadata.State)

		plugin.state = append(plugin.state, stateType)
	}

	for _, action := range builder.actions {
		actionType := action.target
		actionType.meta = meta.Member(action.name)
		panics.IsTrue(actionType.meta != nil)
		panics.IsTrue(actionType.meta.MemberType() == metadata.Action)

		plugin.actions = append(plugin.actions, actionType)
	}

	for _, configItem := range builder.config {
		configType := configItem.target
		configType.meta = meta.Config(configItem.name)
		panics.IsTrue(configType.meta != nil)

		plugin.config = append(plugin.config, configType)
	}

	return plugin
}
