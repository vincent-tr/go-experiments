package conditions

import (
	"go-experiments/traders/modular/context"
	"go-experiments/traders/modular/formatter"
	"go-experiments/traders/modular/marshal"
)

type Condition interface {
	formatter.Formatter
	Execute(ctx context.TraderContext) bool
	ToJsonSpec() (string, any)
}

func newCondition(
	execute func(ctx context.TraderContext) bool,
	format func() *formatter.FormatterNode,
	toJsonSpec func() (string, any),
) Condition {
	return &condition{
		execute:    execute,
		format:     format,
		toJsonSpec: toJsonSpec,
	}
}

type condition struct {
	execute    func(ctx context.TraderContext) bool
	format     func() *formatter.FormatterNode
	toJsonSpec func() (string, any)
}

func (c *condition) Execute(ctx context.TraderContext) bool {
	return c.execute(ctx)
}

func (c *condition) Format() *formatter.FormatterNode {
	return c.format()
}

func (c *condition) ToJsonSpec() (string, any) {
	return c.toJsonSpec()
}

var jsonParsers = marshal.NewRegistry[Condition]()

func FromJSON(jsonData []byte) (Condition, error) {
	return jsonParsers.FromJSON(jsonData)
}
