package opentimecondition

import (
	"go-experiments/traders/modular/formatter"
	"time"
)

type OpenTimeCondition interface {
	formatter.Formatter
	Execute(timestamp time.Time) bool
}

func And(conditions ...OpenTimeCondition) OpenTimeCondition {
	return &andCondition{conditions: conditions}
}

type andCondition struct {
	conditions []OpenTimeCondition
}

func (a *andCondition) Execute(timestamp time.Time) bool {
	for _, condition := range a.conditions {
		if !condition.Execute(timestamp) {
			return false
		}
	}
	return true
}

func (a *andCondition) Format() *formatter.FormatterNode {
	return formatter.FormatWithChildren("And", a.conditions...)
}

func Or(conditions ...OpenTimeCondition) OpenTimeCondition {
	return &orCondition{conditions: conditions}
}

type orCondition struct {
	conditions []OpenTimeCondition
}

func (o *orCondition) Execute(timestamp time.Time) bool {
	for _, condition := range o.conditions {
		if condition.Execute(timestamp) {
			return true
		}
	}
	return false
}

func (o *orCondition) Format() *formatter.FormatterNode {
	return formatter.FormatWithChildren("Or", o.conditions...)
}
