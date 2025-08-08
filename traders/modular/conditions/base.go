package conditions

import (
	"go-experiments/traders/modular/context"
	"go-experiments/traders/modular/formatter"
	"go-experiments/traders/modular/json"
)

type Condition interface {
	formatter.Formatter
	Execute(ctx context.TraderContext) bool
}

func newCondition(
	execute func(ctx context.TraderContext) bool,
	format func() *formatter.FormatterNode,
) Condition {
	return &condition{
		execute: execute,
		format:  format,
	}
}

type condition struct {
	execute func(ctx context.TraderContext) bool
	format  func() *formatter.FormatterNode
}

func (c *condition) Execute(ctx context.TraderContext) bool {
	return c.execute(ctx)
}

func (c *condition) Format() *formatter.FormatterNode {
	return c.format()
}

var jsonParsers = json.NewRegistry[Condition]()

func FromJSON(jsonData []byte) (Condition, error) {
	return jsonParsers.FromJSON(jsonData)
}

func ToJSON(condition Condition) ([]byte, error) {
	panic("ToJSON not implemented for conditions")
}
