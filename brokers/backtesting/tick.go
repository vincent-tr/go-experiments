package backtesting

import "time"

// Tick represents one row of tick data
type tick struct {
	Timestamp time.Time
	Bid       float64
	Ask       float64
}

func (t *tick) Price() float64 {
	// For simplicity, we return the average of bid and ask as the price.
	// In a real implementation, you might want to use bid or ask based on your strategy.
	return (t.Bid + t.Ask) / 2
}
