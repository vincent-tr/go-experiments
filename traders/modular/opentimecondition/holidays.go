package opentimecondition

import (
	"go-experiments/common"
	"go-experiments/traders/modular/formatter"
	"time"
)

func ExcludeUKHolidays() OpenTimeCondition {
	return &excludeCallbackCondition{
		name:     "ExcludeUKHolidays",
		callback: common.IsUKHoliday,
	}
}

func ExcludeUSHolidays() OpenTimeCondition {
	return &excludeCallbackCondition{
		name:     "ExcludeUSHolidays",
		callback: common.IsUSHoliday,
	}
}

type excludeCallbackCondition struct {
	name     string
	callback func(t time.Time) bool
}

func (e *excludeCallbackCondition) Execute(timestamp time.Time) bool {
	return !e.callback(timestamp)
}

func (e *excludeCallbackCondition) Format() *formatter.FormatterNode {
	return formatter.Format(e.name)
}
