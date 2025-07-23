package condition

import (
	"fmt"
	"go-experiments/traders/modular/context"
	"go-experiments/traders/modular/formatter"
	"strings"
	"time"
)

func Weekday(weekdays ...time.Weekday) Condition {
	return &weekdayCondition{
		weekdays: weekdays,
	}
}

type weekdayCondition struct {
	weekdays []time.Weekday
}

func (w *weekdayCondition) Execute(ctx context.TraderContext) bool {
	for _, day := range w.weekdays {
		if ctx.Timestamp().Weekday() == day {
			return true
		}
	}
	return false
}
func (w *weekdayCondition) Format() *formatter.FormatterNode {
	weekdays := make([]string, len(w.weekdays))
	for i, day := range w.weekdays {
		weekdays[i] = day.String()
	}

	return formatter.Format(fmt.Sprintf("Weekday: %s", strings.Join(weekdays, ", ")))
}
