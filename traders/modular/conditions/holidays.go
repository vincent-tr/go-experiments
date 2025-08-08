package conditions

import (
	"encoding/json"
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
		func() (string, any) {
			return "excludeUKHolidays", nil
		},
	)
}

func init() {
	jsonParsers.RegisterParser("excludeUKHolidays", func(arg json.RawMessage) (Condition, error) {
		return ExcludeUKHolidays(), nil
	})
}

func ExcludeUSHolidays() Condition {
	return newCondition(
		func(ctx context.TraderContext) bool {
			return !common.IsUSHoliday(ctx.Timestamp())
		},
		func() *formatter.FormatterNode {
			return formatter.Format("ExcludeUSHolidays")
		},
		func() (string, any) {
			return "excludeUSHolidays", nil
		},
	)
}

func init() {
	jsonParsers.RegisterParser("excludeUSHolidays", func(arg json.RawMessage) (Condition, error) {
		return ExcludeUSHolidays(), nil
	})
}
