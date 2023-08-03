package metadata

type ConfigType string

const (
	String  ConfigType = "string"
	Bool               = "bool"
	Integer            = "integer"
	Float              = "float"
)

type ConfigItem struct {
	name        string
	description string
	valueType   ConfigType
}

func (this *ConfigItem) Name() string {
	return this.name
}

func (this *ConfigItem) Description() string {
	return this.description
}

func (this *ConfigItem) ValueType() ConfigType {
	return this.valueType
}
