package main

import (
	"context"
	"time"
)

type Trader1 struct {
}

// GetTimeframe implements Trader.
func (t *Trader1) GetTimeframe() time.Duration {
	return time.Minute
}

// Run implements Trader.
func (t *Trader1) Run(ctx context.Context, runner Runner) error {
	for {
		candle, err := runner.WaitNext(ctx)
		if err != nil {
			return err
		}

		if false {
			order := &Order{
				Direction: PositionDirectionLong,
				Price:     42,
				Reason:    "Candle closed higher",
			}
			_, err = runner.PlaceOrder(order)
			if err != nil {
				return err
			}
		}
	}
}

var _ Trader = (*Trader1)(nil)
