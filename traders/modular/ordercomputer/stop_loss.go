package ordercomputer

import (
	"fmt"
	"go-experiments/brokers"
	"go-experiments/traders/modular/context"
	"go-experiments/traders/modular/formatter"

	"github.com/markcheno/go-talib"
)

func StopLossATR(period int, multiplier float64) OrderComputer {
	return newOrderComputer(
		func(ctx context.TraderContext, order *brokers.Order) error {
			history := ctx.HistoricalData()
			atr := talib.Atr(history.GetHighPrices(), history.GetLowPrices(), history.GetClosePrices(), period)

			if len(atr) == 0 {
				return fmt.Errorf("not enough data for ATR calculation")
			}

			lastAtr := atr[len(atr)-1]
			order.StopLoss = lastAtr * multiplier
			return nil
		},
		func() *formatter.FormatterNode {
			return formatter.Format("StopLossATR",
				formatter.Format(fmt.Sprintf("Period: %d", period)),
				formatter.Format(fmt.Sprintf("Multiplier: %.4f", multiplier)),
			)
		},
	)
}

func StopLossPipBuffer(pipBuffer int, lookupPeriod int) OrderComputer {
	return newOrderComputer(
		func(ctx context.TraderContext, order *brokers.Order) error {
			pipSize := 0.0001 // Assuming a pip size of 0.0001 for most currency pairs
			pipDistance := float64(pipBuffer) * pipSize

			switch order.Direction {
			case brokers.PositionDirectionLong:
				// find lowest low in last lookupPeriod minutes
				lowest := ctx.HistoricalData().GetLowest(lookupPeriod)
				order.StopLoss = lowest - pipDistance
				return nil
			case brokers.PositionDirectionShort:
				// find highest high in last lookupPeriod minutes
				highest := ctx.HistoricalData().GetHighest(lookupPeriod)
				order.StopLoss = highest + pipDistance
				return nil
			default:
				return fmt.Errorf("invalid position direction: %s", order.Direction.String())
			}
		},
		func() *formatter.FormatterNode {
			return formatter.Format("StopLossPipBuffer",
				formatter.Format(fmt.Sprintf("Pip Buffer: %d", pipBuffer)),
				formatter.Format(fmt.Sprintf("Lookup Period: %d", lookupPeriod)),
			)
		},
	)
}
