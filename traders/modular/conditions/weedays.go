package conditions

import (
	"encoding/json"
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
		func() (string, any) {
			weekdaysStr := make([]string, len(weekdays))
			for i, day := range weekdays {
				weekdaysStr[i] = day.String()
			}
			return "weekday", map[string]any{
				"weekdays": weekdaysStr,
			}
		},
	)
}

func init() {
	jsonParsers.RegisterParser("weekday", func(arg json.RawMessage) (Condition, error) {
		var weekdays []string

		if err := json.Unmarshal(arg, &weekdays); err != nil {
			return nil, fmt.Errorf("failed to parse weekday condition: %w", err)
		}

		var days []time.Weekday
		for _, dayStr := range weekdays {
			day, err := time.Parse("Monday", dayStr)
			if err != nil {
				return nil, fmt.Errorf("invalid weekday: %s", dayStr)
			}
			days = append(days, day.Weekday())
		}

		return Weekday(days...), nil
	})
}
