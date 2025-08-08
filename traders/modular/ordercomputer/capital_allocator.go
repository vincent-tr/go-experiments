package ordercomputer

import (
	"encoding/json"
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
			return formatter.Format(fmt.Sprintf("CapitalRiskPercent: %.2f%%", riskPerTradePercent))
		},
	)
}

func init() {
	jsonParsers.RegisterParser("capitalRiskPercent", func(arg json.RawMessage) (OrderComputer, error) {
		var riskPerTradePercent float64
		if err := json.Unmarshal(arg, &riskPerTradePercent); err != nil {
			return nil, fmt.Errorf("failed to parse CapitalRiskPercent: %w", err)
		}

		return CapitalRiskPercent(riskPerTradePercent), nil
	})
}

func CapitalFixed(amount float64) OrderComputer {
	return newOrderComputer(
		func(ctx context.TraderContext, order *brokers.Order) error {
			broker := ctx.Broker()
			accountBalance := broker.GetCapital()
			accountRisk := amount

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
			return formatter.Format(fmt.Sprintf("CapitalFixed: %.2f", amount))
		},
	)
}

func init() {
	jsonParsers.RegisterParser("capitalFixed", func(arg json.RawMessage) (OrderComputer, error) {
		var amount float64
		if err := json.Unmarshal(arg, &amount); err != nil {
			return nil, fmt.Errorf("failed to parse CapitalFixed: %w", err)
		}

		return CapitalFixed(amount), nil
	})
}
