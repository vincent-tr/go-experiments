package main

import (
	"context"
	"time"
)

// Trader is an interface that defines a trading strategy.
// A trader is a strategy that can be run on a Runner.
// It defines the methods that a trader must implement to be able to run on a Runner.
type Trader interface {
	// Get the timeframe of the trader.
	// This is the duration of each candle that the trader will use.
	// For example, if the trader is using 1 minute candles, this will return time.Minute.
	GetTimeframe() time.Duration

	// Run the trader with the given runner and context.
	Run(ctx context.Context, broker Broker) error
}
