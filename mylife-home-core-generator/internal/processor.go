package internal

import (
	annotation "github.com/YReshetko/go-annotation/pkg"
)

func init() {
	annotation.Register[SomeAnnotationStructure](&SomeAnnotationProcessor{})
}

var _ annotation.AnnotationProcessor = (*SomeAnnotationProcessor)(nil)

type SomeAnnotationProcessor struct{}

func (p *SomeAnnotationProcessor) Process(_ annotation.Node) error {
	// Single node processing
}

func (p *SomeAnnotationProcessor) Output() map[string][]byte {
	// Prepare processing results and return:
	// map.key (string) - absolute file path (processed dir can be taken from annotation.Node: node.Dir())
	// map.value ([]byte) - resulting file data (for example .go file)
}

func (p *SomeAnnotationProcessor) Version() string {
	return "0.0.1" // any string, that represents the processor version
}

func (p *SomeAnnotationProcessor) Name() string {
	// Any string that represents the processor name,
	//if the processor handles a single annotation it could be the annotation name
	return "SomeAnnotationStructure"
}
