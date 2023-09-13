package internal

import (
	"fmt"
	"path"

	annotation "github.com/YReshetko/go-annotation/pkg"
)

type Generator struct {
	outputPath string
}

func MakeGenerator(node annotation.Node, outputPath string) *Generator {
	return &Generator{outputPath}
}

func (generator *Generator) EnsureOutputPath(node annotation.Node) {
	if len(generator.outputPath) == 0 {
		filename := node.Meta().FileName()
		filename = filename[:len(filename)-2] + "annotations-generated.go"
		generator.outputPath = path.Join(node.Meta().Dir(), filename)
		// TODO: per input file (can have several plugins in several files)
		fmt.Println(generator.outputPath)
	}
}

func (generator *Generator) ProcessPluginAnnotations(node annotation.Node, annotations []Plugin) {
	fmt.Println("pluginAnnotations")
	fmt.Println(annotations)
}

func (generator *Generator) ProcessStateAnnotations(node annotation.Node, annotations []State) {
	fmt.Println("stateAnnotations")
	fmt.Println(annotations)
}

func (generator *Generator) ProcessActionAnnotations(node annotation.Node, annotations []Action) {
	fmt.Println("actionAnnotations")
	fmt.Println(annotations)
}

func (generator *Generator) ProcessConfigAnnotations(node annotation.Node, annotations []Config) {
	fmt.Println("configAnnotations")
	fmt.Println(annotations)
}

func (generator *Generator) Output() []byte {
	return make([]byte, 0)
}
