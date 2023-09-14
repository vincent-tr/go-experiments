package registry

type Registry interface {
	AddPlugin(pluginType *PluginType)
}

var registry Registry

func RegisterPlugin(pluginType *PluginType) {
	registry.AddPlugin(pluginType)
}

func SetRegistry(value Registry) {
	registry = value
}
