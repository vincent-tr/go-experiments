package internal

import (
	"fmt"
	"path"

	annotation "github.com/YReshetko/go-annotation/pkg"
)

func init() {
	processor := &PluginAnnotationProcessor{}

	annotation.Register[Plugin](processor)
	annotation.Register[Config](processor)
	annotation.Register[State](processor)
	annotation.Register[Action](processor)
}

var _ annotation.AnnotationProcessor = (*PluginAnnotationProcessor)(nil)

type PluginAnnotationProcessor struct {
	outputPath string
}

func (this *PluginAnnotationProcessor) Process(node annotation.Node) error {
	if len(this.outputPath) == 0 {
		filename := node.Meta().FileName()
		filename = filename[:len(filename)-2] + "annotations-generated.go"
		// TODO: per input file
		this.outputPath = path.Join(node.Meta().Dir(), filename)
		fmt.Println(this.outputPath)
	}

	pluginAnnotations := annotation.FindAnnotations[Plugin](node.Annotations())
	fmt.Println(pluginAnnotations)

	configAnnotations := annotation.FindAnnotations[Config](node.Annotations())
	fmt.Println(configAnnotations)

	stateAnnotations := annotation.FindAnnotations[State](node.Annotations())
	fmt.Println(stateAnnotations)

	actionAnnotations := annotation.FindAnnotations[Action](node.Annotations())
	fmt.Println(actionAnnotations)

	return nil
}

func (this *PluginAnnotationProcessor) Output() map[string][]byte {
	// Prepare processing results and return:
	// map.key (string) - absolute file path (processed dir can be taken from annotation.Node: node.Dir())
	// map.value ([]byte) - resulting file data (for example .go file)
	return make(map[string][]byte)
}

func (this *PluginAnnotationProcessor) Version() string {
	return "0.0.1"
}

func (this *PluginAnnotationProcessor) Name() string {
	return "MylifeHomePluginAnnotationProcessor"
}
