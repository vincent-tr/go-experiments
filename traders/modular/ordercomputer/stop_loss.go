package ordercomputer

import (
	"fmt"
	"go-experiments/brokers"
	"go-experiments/traders/modular/context"
	"go-experiments/traders/modular/formatter"
	"go-experiments/traders/modular/indicators"
)

func StopLossATR(atr indicators.Indicator, multiplier float64) OrderComputer {
	return newOrderComputer(
		func(ctx context.TraderContext, order *brokers.Order) error {
			atr := atr.Values(ctx)

			if len(atr) == 0 {
				return fmt.Errorf("not enough data for ATR calculation")
			}

			currAtr := atr[len(atr)-1]
			pipDistance := currAtr * multiplier
			entryPrice := ctx.EntryPrice()

			switch order.Direction {
			case brokers.PositionDirectionLong:
				order.StopLoss = entryPrice - pipDistance
				return nil

			case brokers.PositionDirectionShort:
				order.StopLoss = entryPrice + pipDistance
				return nil

			default:
				panic("invalid position type")
			}
		},
		func() *formatter.FormatterNode {
			return formatter.Format("StopLossATR",
				formatter.FormatWithChildren("ATR", atr),
				formatter.Format(fmt.Sprintf("Multiplier: %.4f", multiplier)),
			)
		},
	)
}

const pipSize = 0.0001

func StopLossPipBuffer(pipBuffer int, lookupPeriod int) OrderComputer {
	pipDistance := float64(pipBuffer) * pipSize

	return newOrderComputer(
		func(ctx context.TraderContext, order *brokers.Order) error {

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
