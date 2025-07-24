package condition

import (
	"go-experiments/traders/modular/context"
	"go-experiments/traders/modular/formatter"
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

func HistoryComplete() Condition {
	return newCondition(
		func(ctx context.TraderContext) bool {
			return ctx.HistoricalData().IsComplete()
		},
		func() *formatter.FormatterNode {
			return formatter.Format("HistoryComplete")
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
