package condition

import (
	"go-experiments/common"
	"go-experiments/traders/modular/context"
	"go-experiments/traders/modular/formatter"
)

func ExcludeUKHolidays() Condition {
	return newCondition(
		func(ctx context.TraderContext) bool {
			return !common.IsUKHoliday(ctx.Timestamp())
		},
		func() *formatter.FormatterNode {
			return formatter.Format("ExcludeUKHolidays")
		},
	)
}

func ExcludeUSHolidays() Condition {
	return newCondition(
		func(ctx context.TraderContext) bool {
			return !common.IsUSHoliday(ctx.Timestamp())
		},
		func() *formatter.FormatterNode {
			return formatter.Format("ExcludeUSHolidays")
		},
	)
}
