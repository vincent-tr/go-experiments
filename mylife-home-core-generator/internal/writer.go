package internal

import (
	"fmt"
	"mylife-home-core-common/metadata"
	"strings"
)

type Writer struct {
	builder *strings.Builder
}

func MakeWrite(packageName string) *Writer {
	writer := &Writer{
		builder: &strings.Builder{},
	}

	writer.appendLinef(`package %s`, packageName)
	writer.appendLine(``)
	writer.appendLine(`import (`)
	writer.appendLine(`	"mylife-home-core-common/metadata"`)
	writer.appendLine(`	"mylife-home-core-common/registry"`)
	writer.appendLine(`)`)
	writer.appendLine(``)

	return writer
}

func (writer *Writer) appendBlock(str string) {
	writer.builder.WriteString(str)
}

func (writer *Writer) appendLine(str string) {
	writer.appendBlock(str + "\n")
}

func (writer *Writer) appendLinef(format string, a ...any) {
	writer.appendLine(fmt.Sprintf(format, a...))
}

func (writer *Writer) BeginPlugin(pluginType string, name string, description string, usage metadata.PluginUsage) {
	writer.appendLine(`func init() {`)
	writer.appendLinef(`	builder := registry.MakePluginTypeBuilder[%s]("%s", "%s", metadata.%s)`, pluginType, name, description, string(usage))
}

func (writer *Writer) AddState(fieldName string, name string, description string, valueType metadata.Type) {
	writer.appendLinef(`	builder.AddState("%s", "%s", "%s", %s)`, fieldName, name, description, renderType(valueType))
}

func (writer *Writer) AddAction(methodName string, name string, description string, valueType metadata.Type) {
	writer.appendLinef(`	builder.AddAction("%s", "%s", "%s", %s)`, methodName, name, description, renderType(valueType))

}

func (writer *Writer) AddConfig(fieldName string, name string, description string, valueType metadata.ConfigType) {
	writer.appendLinef(`	builder.AddConfig("%s", "%s", "%s", %s)`, fieldName, name, description, string(valueType))
}

func (writer *Writer) EndPlugin() {
	writer.appendLine(`	registry.RegisterPlugin(builder.Build())`)
	writer.appendLine(`}`)
	writer.appendLine(``)
}

func (writer *Writer) Content() []byte {
	return []byte(writer.builder.String())
}

func renderType(typ metadata.Type) string {
	switch typed := typ.(type) {

	case *metadata.RangeType:
		return fmt.Sprintf(`MakeTypeRange(%d, %d)`, typed.Min(), typed.Max())

	case *metadata.TextType:
		return `MakeTypeText()`

	case *metadata.FloatType:
		return `MakeTypeFloat()`

	case *metadata.BoolType:
		return `MakeTypeBool()`

	case *metadata.EnumType:
		builder := strings.Builder{}
		builder.WriteString(`MakeTypeEnum(`)
		for index := 0; index < typed.NumValues(); index += 1 {
			if index > 0 {
				builder.WriteString(`, `)
			}

			builder.WriteString(`"` + typed.Value(index) + `"`)
		}
		builder.WriteString(`)`)
		return builder.String()

	case *metadata.ComplexType:
		return `MakeTypeComplex()`

	default:
		return "???"
	}
}

/*
package plugin

import (
	"mylife-home-core-common/metadata"
	"mylife-home-core-common/registry"
)

func init() {
	builder := registry.MakePluginTypeBuilder[ValueBinary]("name", "desc", metadata.Logic)
	builder.AddState("Value", "value", "desc", metadata.MakeTypeBool())
	builder.AddAction("SetValue", "setValue", "desc", metadata.MakeTypeBool())
	builder.AddConfig("InitialValue", "initialValue", "desc", metadata.Bool)
	registry.RegisterPlugin(builder.Build())
}
*/
