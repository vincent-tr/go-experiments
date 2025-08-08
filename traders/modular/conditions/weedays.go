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
				weekdaysStr[i] = strings.ToLower(day.String())
			}
			return "weekday", weekdaysStr
		},
	)
}

func init() {
	parser := make(map[string]time.Weekday)

	for i := 0; i < 7; i++ {
		day := time.Weekday(i)
		parser[strings.ToLower(day.String())] = day
	}

	jsonParsers.RegisterParser("weekday", func(arg json.RawMessage) (Condition, error) {
		var weekdays []string

		if err := json.Unmarshal(arg, &weekdays); err != nil {
			return nil, fmt.Errorf("failed to parse weekday condition: %w", err)
		}

		var days []time.Weekday
		for _, dayStr := range weekdays {
			day, ok := parser[strings.ToLower(dayStr)]
			if !ok {
				return nil, fmt.Errorf("invalid weekday: %s", dayStr)
			}
			days = append(days, day)
		}

		return Weekday(days...), nil
	})
}
