package metadata

type serializerImpl struct{}

var Seralizer = serializerImpl{}

func (e *serializerImpl) SerializeComponent(component *Component) any {
	panic("TODO")
}

func (e *serializerImpl) DeserializeComponent(data any) *Component {
	panic("TODO")
}

func (e *serializerImpl) SerializePlugin(plugin *Plugin) any {
	panic("TODO")
}

func (e *serializerImpl) DeserializePlugin(data any) *Plugin {
	panic("TODO")
}
