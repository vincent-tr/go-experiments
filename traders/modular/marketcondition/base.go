package marketcondition

import (
	"go-experiments/traders/modular/formatter"
	"go-experiments/traders/tools"
)

type MarketCondition interface {
	formatter.Formatter
	Execute(history *tools.History) bool
}

func And(conditions ...MarketCondition) MarketCondition {
	return &andCondition{conditions: conditions}
}

type andCondition struct {
	conditions []MarketCondition
}

func (a *andCondition) Execute(history *tools.History) bool {
	for _, condition := range a.conditions {
		if !condition.Execute(history) {
			return false
		}
	}
	return true
}

func (a *andCondition) Format() *formatter.FormatterNode {
	return formatter.FormatWithChildren("And", a.conditions...)
}

type orCondition struct {
	conditions []MarketCondition
}

func (o *orCondition) Execute(history *tools.History) bool {
	for _, condition := range o.conditions {
		if condition.Execute(history) {
			return true
		}
	}
	return false
}

func (o *orCondition) Format() *formatter.FormatterNode {
	return formatter.FormatWithChildren("Or", o.conditions...)
}

func HistoryComplete() MarketCondition {
	return &historyCompleteCondition{}
}

type historyCompleteCondition struct{}

func (h *historyCompleteCondition) Execute(history *tools.History) bool {
	return history.IsComplete()
}

func (h *historyCompleteCondition) Format() *formatter.FormatterNode {
	return formatter.Format("HistoryComplete")
}
