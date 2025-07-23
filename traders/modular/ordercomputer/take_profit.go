package ordercomputer

import (
	"fmt"
	"go-experiments/brokers"
	"go-experiments/traders/modular/context"
	"go-experiments/traders/modular/formatter"
)

func TakeProfitRatio(ratio float64) OrderComputer {
	return &takeProfitRatio{
		ratio: ratio,
	}
}

type takeProfitRatio struct {
	ratio float64
}

func (oc *takeProfitRatio) Compute(ctx context.TraderContext, order *brokers.Order) error {
	if order.StopLoss == 0 {
		return fmt.Errorf("stop loss must be set before calculating take profit")
	}

	entryPrice := ctx.EntryPrice()

	switch order.Direction {
	case brokers.PositionDirectionLong:
		risk := entryPrice - order.StopLoss
		if risk <= 0 {
			return fmt.Errorf("invalid stoploss for long position: entryPrice=%.5f, stopLoss=%.5f", entryPrice, order.StopLoss)
		}
		order.TakeProfit = entryPrice + oc.ratio*risk
		return nil

	case brokers.PositionDirectionShort:
		risk := order.StopLoss - entryPrice
		if risk <= 0 {
			return fmt.Errorf("invalid stoploss for short position: entryPrice=%.5f, stopLoss=%.5f", entryPrice, order.StopLoss)
		}
		order.TakeProfit = entryPrice - oc.ratio*risk
		return nil

	default:
		return fmt.Errorf("invalid position direction: %s", order.Direction.String())
	}
}

func (oc *takeProfitRatio) Format() *formatter.FormatterNode {
	return formatter.Format("TakeProfitRatio",
		formatter.Format(fmt.Sprintf("Ratio: %.4f", oc.ratio)),
	)
}
