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

func (config *ConfigItem) Name() string {
	return config.name
}

func (config *ConfigItem) Description() string {
	return config.description
}

func (config *ConfigItem) ValueType() ConfigType {
	return config.valueType
}
