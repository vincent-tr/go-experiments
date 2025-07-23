package opentimecondition

import (
	"fmt"
	"go-experiments/traders/modular/formatter"
	"strings"
	"time"
)

func Weekday(weekdays ...time.Weekday) OpenTimeCondition {
	return &weekdayCondition{
		weekdays: weekdays,
	}
}

type weekdayCondition struct {
	weekdays []time.Weekday
}

func (w *weekdayCondition) Execute(timestamp time.Time) bool {
	for _, day := range w.weekdays {
		if timestamp.Weekday() == day {
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
