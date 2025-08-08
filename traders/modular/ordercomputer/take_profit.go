package ordercomputer

import (
	"encoding/json"
	"fmt"
	"go-experiments/brokers"
	"go-experiments/traders/modular/context"
	"go-experiments/traders/modular/formatter"
)

func TakeProfitRatio(ratio float64) OrderComputer {
	return newOrderComputer(
		func(ctx context.TraderContext, order *brokers.Order) error {
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
				order.TakeProfit = entryPrice + ratio*risk
				return nil

			case brokers.PositionDirectionShort:
				risk := order.StopLoss - entryPrice
				if risk <= 0 {
					return fmt.Errorf("invalid stoploss for short position: entryPrice=%.5f, stopLoss=%.5f", entryPrice, order.StopLoss)
				}
				order.TakeProfit = entryPrice - ratio*risk
				return nil

			default:
				return fmt.Errorf("invalid position direction: %s", order.Direction.String())
			}
		},
		func() *formatter.FormatterNode {
			return formatter.Format(fmt.Sprintf("TakeProfitRatio: %.4f", ratio))
		},
		func() (string, any) {
			return "takeProfitRatio", ratio
		},
	)
}

func init() {
	jsonParsers.RegisterParser("takeProfitRatio", func(arg json.RawMessage) (OrderComputer, error) {
		var ratio float64
		if err := json.Unmarshal(arg, &ratio); err != nil {
			return nil, fmt.Errorf("failed to parse TakeProfitRatio: %w", err)
		}

		return TakeProfitRatio(ratio), nil
	})
}
