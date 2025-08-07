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

func And(conditions ...Condition) Condition {
	return newCondition(
		func(ctx context.TraderContext) bool {
			for _, condition := range conditions {
				if !condition.Execute(ctx) {
					return false
				}
			}
			return true
		},
		func() *formatter.FormatterNode {
			return formatter.FormatWithChildren("And", conditions...)
		},
	)
}

func Or(conditions ...Condition) Condition {
	return newCondition(
		func(ctx context.TraderContext) bool {
			for _, condition := range conditions {
				if condition.Execute(ctx) {
					return true
				}
			}
			return false
		},
		func() *formatter.FormatterNode {
			return formatter.FormatWithChildren("Or", conditions...)
		},
	)
}

func HistoryUsable() Condition {
	return newCondition(
		func(ctx context.TraderContext) bool {
			return ctx.HistoricalData().IsUsable()
		},
		func() *formatter.FormatterNode {
			return formatter.Format("HistoryUsable")
		},
	)
}

func NoOpenPositions() Condition {
	return newCondition(
		func(ctx context.TraderContext) bool {
			return len(ctx.OpenPositions()) == 0
		},
		func() *formatter.FormatterNode {
			return formatter.Format("NoOpenPositions")
		},
	)
}
