package condition

import (
	"go-experiments/common"
	"go-experiments/traders/modular/context"
	"go-experiments/traders/modular/formatter"
	"time"
)

func ExcludeUKHolidays() Condition {
	return &excludeCallbackCondition{
		name:     "ExcludeUKHolidays",
		callback: common.IsUKHoliday,
	}
}

func ExcludeUSHolidays() Condition {
	return &excludeCallbackCondition{
		name:     "ExcludeUSHolidays",
		callback: common.IsUSHoliday,
	}
}

type excludeCallbackCondition struct {
	name     string
	callback func(t time.Time) bool
}

func (e *excludeCallbackCondition) Execute(ctx context.TraderContext) bool {
	return !e.callback(ctx.Timestamp())
}

func (e *excludeCallbackCondition) Format() *formatter.FormatterNode {
	return formatter.Format(e.name)
}
