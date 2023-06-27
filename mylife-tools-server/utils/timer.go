package utils

import "time"

type Timer struct {
	begin time.Time
}

func NewTimer() Timer {
	return Timer{begin: time.Now()}
}

func (t *Timer) ElapsedMs() float64 {
	return time.Since(t.begin).Seconds() * 1000
}
