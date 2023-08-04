package internal

import (
	"fmt"
	annotation "github.com/YReshetko/go-annotation/pkg"
)

func init() {
	annotation.Register[Component](&ComponentAnnotationProcessor{})
}

var _ annotation.AnnotationProcessor = (*ComponentAnnotationProcessor)(nil)

type ComponentAnnotationProcessor struct {
	root string
}

func (this *ComponentAnnotationProcessor) Process(node annotation.Node) error {
	// Single node processing
	if len(this.root) == 0 {
		this.root = node.Meta().Root()
	}
	ans := annotation.FindAnnotations[Component](node.Annotations())
	fmt.Println(ans)

	return nil
}

func (this *ComponentAnnotationProcessor) Output() map[string][]byte {
	// Prepare processing results and return:
	// map.key (string) - absolute file path (processed dir can be taken from annotation.Node: node.Dir())
	// map.value ([]byte) - resulting file data (for example .go file)
	return make(map[string][]byte)
}

func (this *ComponentAnnotationProcessor) Version() string {
	return "0.0.1" // any string, that represents the processor version
}

func (this *ComponentAnnotationProcessor) Name() string {
	// Any string that represents the processor name,
	//if the processor handles a single annotation it could be the annotation name
	return "ComponentAnnotationProcessor"
}
