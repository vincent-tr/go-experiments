package condition

import (
	"go-experiments/traders/modular/context"
	"go-experiments/traders/modular/formatter"
)

type Condition interface {
	formatter.Formatter
	Execute(ctx context.TraderContext) bool
}

func And(conditions ...Condition) Condition {
	return &andCondition{conditions: conditions}
}

type andCondition struct {
	conditions []Condition
}

func (a *andCondition) Execute(ctx context.TraderContext) bool {
	for _, condition := range a.conditions {
		if !condition.Execute(ctx) {
			return false
		}
	}
	return true
}

func (a *andCondition) Format() *formatter.FormatterNode {
	return formatter.FormatWithChildren("And", a.conditions...)
}

func Or(conditions ...Condition) Condition {
	return &orCondition{conditions: conditions}
}

type orCondition struct {
	conditions []Condition
}

func (o *orCondition) Execute(ctx context.TraderContext) bool {
	for _, condition := range o.conditions {
		if condition.Execute(ctx) {
			return true
		}
	}
	return false
}

func (o *orCondition) Format() *formatter.FormatterNode {
	return formatter.FormatWithChildren("Or", o.conditions...)
}

func HistoryComplete() Condition {
	return &historyCompleteCondition{}
}

type historyCompleteCondition struct{}

func (h *historyCompleteCondition) Execute(ctx context.TraderContext) bool {
	return ctx.HistoricalData().IsComplete()
}

func (h *historyCompleteCondition) Format() *formatter.FormatterNode {
	return formatter.Format("HistoryComplete")
}

func OnePositionAtATime() Condition {
	return &onePositionAtATimeCondition{}
}

type onePositionAtATimeCondition struct{}

func (o *onePositionAtATimeCondition) Execute(ctx context.TraderContext) bool {
	return len(ctx.OpenPositions()) == 0
}

func (o *onePositionAtATimeCondition) Format() *formatter.FormatterNode {
	return formatter.Format("OnePositionAtATime")
}
