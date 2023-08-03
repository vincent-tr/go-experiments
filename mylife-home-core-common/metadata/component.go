package metadata

type PluginUsage string

const (
	Sensor   PluginUsage = "sensor"
	Actuator             = "actuator"
	Logic                = "logic"
	Ui                   = "ui"
)

type Component struct {
	name        string
	description string
	usage       PluginUsage
	config      map[string]*ConfigItem
	members     map[string]*Member
}

func (this *Component) Name() string {
	return this.name
}

func (this *Component) Description() string {
	return this.description
}

func (this *Component) Usage() PluginUsage {
	return this.usage
}

func (this *Component) ConfigNames() string {
	return maps.Keys(this.config)
}

func (this *Component) Config(name string) *ConfigItem {
	return this.config[name]
}

func (this *Component) MemberNames() string {
	return maps.Keys(this.members)
}

func (this *Component) Member(name string) *Member {
	return this.members[name]
}
