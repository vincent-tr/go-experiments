package backtesting

import "time"

// Tick represents one row of tick data
type tick struct {
	Timestamp time.Time
	Bid       float64
	Ask       float64
}
