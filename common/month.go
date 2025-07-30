package common

import (
	"fmt"
	"time"
)

type Month struct {
	year  int
	month int
}

func NewMonth(year int, month int) Month {
	return Month{year: year, month: month}
}

func FromDate(t time.Time) Month {
	return Month{year: t.Year(), month: int(t.Month())}
}

func (m Month) Year() int {
	return m.year
}

func (m Month) Month() int {
	return m.month
}

func (m Month) String() string {
	return fmt.Sprintf("%04d-%02d", m.year, m.month)
}

func (m Month) FirstDay() time.Time {
	return time.Date(m.year, time.Month(m.month), 1, 0, 0, 0, 0, time.UTC)
}

func (m Month) LastDay() time.Time {
	return time.Date(m.year, time.Month(m.month), 1, 0, 0, 0, 0, time.UTC).AddDate(0, 1, -1)
}
