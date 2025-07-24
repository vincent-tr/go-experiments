package ordercomputer

import (
	"fmt"
	"go-experiments/brokers"
	"go-experiments/traders/modular/context"
	"go-experiments/traders/modular/formatter"
	"math"
)

func CapitalRiskPercent(riskPerTradePercent float64) OrderComputer {
	return newOrderComputer(
		func(ctx context.TraderContext, order *brokers.Order) error {
			broker := ctx.Broker()
			accountBalance := broker.GetCapital()
			accountRisk := accountBalance * (riskPerTradePercent / 100)

			entryPrice := ctx.EntryPrice()
			priceDiff := math.Abs(entryPrice - order.StopLoss)
			if priceDiff <= 0 {
				return fmt.Errorf("invalid stop loss price: entryPrice=%.5f, stopLoss=%.5f", entryPrice, order.StopLoss)
			}

			lotSize := float64(broker.GetLotSize())
			riskPerLot := lotSize * priceDiff
			positionSize := accountRisk / riskPerLot

			// Ensure position size doesn't exceed account balance
			// Total value = positionSize * lotSize * entryPrice
			maxPositionSize := accountBalance*broker.GetLeverage()/(lotSize*entryPrice) - 1
			maxPositionSize -= 1 // Avoid rounding issues
			if positionSize > maxPositionSize {
				positionSize = maxPositionSize
			}

			order.Quantity = int(math.Floor(positionSize))
			return nil
		},
		func() *formatter.FormatterNode {
			return formatter.Format("CapitalRiskPercent",
				formatter.Format(fmt.Sprintf("RiskPerTradePercent: %.2f%%", riskPerTradePercent)),
			)
		},
	)
}
