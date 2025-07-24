package condition

import (
	"fmt"
	"go-experiments/traders/modular/context"
	"go-experiments/traders/modular/formatter"
	"strings"
	"time"
)

func Weekday(weekdays ...time.Weekday) Condition {
	return newCondition(
		func(ctx context.TraderContext) bool {
			for _, day := range weekdays {
				if ctx.Timestamp().Weekday() == day {
					return true
				}
			}
			return false
		},
		func() *formatter.FormatterNode {
			weekdaysStr := make([]string, len(weekdays))
			for i, day := range weekdays {
				weekdaysStr[i] = day.String()
			}
			return formatter.Format(fmt.Sprintf("Weekday: %s", strings.Join(weekdaysStr, ", ")))
		},
	)
}
