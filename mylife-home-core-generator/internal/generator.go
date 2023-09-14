package internal

import (
	"fmt"
	"go/ast"

	annotation "github.com/YReshetko/go-annotation/pkg"
	"github.com/gookit/goutil/errorx/panics"
)

type Generator struct {
	outputPath  string
	packageName string

	plugins []*PluginData
	states  []*StateData
	actions []*ActionData
	configs []*ConfigData
}

type PluginData struct {
	typeName string
	typ      *ast.TypeSpec
	ann      *Plugin

	states  []*StateData
	actions []*ActionData
	configs []*ConfigData
}

type StateData struct {
	typeName string
	field    *ast.Field
	ann      *State
}

type ActionData struct {
	typeName string
	fn       *ast.FuncDecl
	ann      *Action
}

type ConfigData struct {
	typeName string
	field    *ast.Field
	ann      *Config
}

func MakeGenerator(node annotation.Node, outputPath string) *Generator {
	return &Generator{
		outputPath:  outputPath,
		packageName: node.Meta().PackageName(),
		plugins:     make([]*PluginData, 0),
		states:      make([]*StateData, 0),
		actions:     make([]*ActionData, 0),
		configs:     make([]*ConfigData, 0),
	}
}

func (generator *Generator) ProcessPluginAnnotation(node annotation.Node, pluginAnnotation *Plugin) {
	typ, ok := annotation.CastNode[*ast.TypeSpec](node)
	panics.IsTrue(ok)
	panics.IsTrue(typ.TypeParams == nil || len(typ.TypeParams.List) == 0)

	generator.plugins = append(generator.plugins, &PluginData{
		typeName: typ.Name.Name,
		typ:      typ,
		ann:      pluginAnnotation,
	})
}

func (generator *Generator) ProcessStateAnnotation(node annotation.Node, ann *State) {
	field, ok := annotation.CastNode[*ast.Field](node)
	panics.IsTrue(ok)

	typ, ok := annotation.ParentType[*ast.TypeSpec](node)
	panics.IsTrue(ok)
	astType, ok := annotation.CastNode[*ast.TypeSpec](typ)
	panics.IsTrue(ok)

	generator.states = append(generator.states, &StateData{
		typeName: astType.Name.Name,
		field:    field,
		ann:      ann,
	})

	fmt.Println("stateAnnotation")
	fmt.Println(ann)
	fmt.Printf("%s %s\n", field.Names[0], field.Type)
}

func (generator *Generator) ProcessActionAnnotation(node annotation.Node, ann *Action) {
	fn, ok := annotation.CastNode[*ast.FuncDecl](node)
	panics.IsTrue(ok)

	astPtr, ok := fn.Recv.List[0].Type.(*ast.StarExpr)
	panics.IsTrue(ok) // Else the function receiver is not a pointer

	generator.actions = append(generator.actions, &ActionData{
		typeName: astPtr.X.(*ast.Ident).Name,
		fn:       fn,
		ann:      ann,
	})

	fmt.Println("actionAnnotation")
	fmt.Println(ann)
	fmt.Printf("%+v\n", fn.Recv.List[0].Type.(*ast.StarExpr).X.(*ast.Ident).Name)
}

func (generator *Generator) ProcessConfigAnnotation(node annotation.Node, ann *Config) {
	field, ok := annotation.CastNode[*ast.Field](node)
	panics.IsTrue(ok)

	typ, ok := annotation.ParentType[*ast.TypeSpec](node)
	panics.IsTrue(ok)
	astType, ok := annotation.CastNode[*ast.TypeSpec](typ)
	panics.IsTrue(ok)

	generator.configs = append(generator.configs, &ConfigData{
		typeName: astType.Name.Name,
		field:    field,
		ann:      ann,
	})

	fmt.Println("configAnnotation")
	fmt.Println(ann)
	fmt.Printf("%s %s\n", field.Names[0], field.Type)
}

func (generator *Generator) Output() []byte {
	generator.associate()

	writer := MakeWrite(generator.packageName)

	return writer.Content()
}

func (generator *Generator) associate() {
	m := make(map[string]*PluginData)

	for _, plugin := range generator.plugins {
		plugin.states = make([]*StateData, 0)
		plugin.actions = make([]*ActionData, 0)
		plugin.configs = make([]*ConfigData, 0)

		m[plugin.typeName] = plugin
	}

	for _, state := range generator.states {
		plugin, ok := m[state.typeName]
		panics.IsTrue(ok)

		plugin.states = append(plugin.states, state)
	}

	for _, action := range generator.actions {
		plugin, ok := m[action.typeName]
		panics.IsTrue(ok)

		plugin.actions = append(plugin.actions, action)
	}

	for _, config := range generator.configs {
		plugin, ok := m[config.typeName]
		panics.IsTrue(ok)

		plugin.configs = append(plugin.configs, config)
	}
}

func (generator *Generator) enrich() {

}
