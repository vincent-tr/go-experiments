package ordercomputer

import (
	"fmt"
	"go-experiments/brokers"
	"go-experiments/traders/modular/formatter"
	"go-experiments/traders/tools"
	"math"
)

func CapitalAllocatorFixedRisk(riskPerTradePercent float64) OrderComputer {
	return &capitalAllocatorFixedRisk{
		riskPerTradePercent: riskPerTradePercent,
	}
}

type capitalAllocatorFixedRisk struct {
	riskPerTradePercent float64
}

func (oc *capitalAllocatorFixedRisk) Compute(broker brokers.Broker, history *tools.History, order *brokers.Order) error {
	accountBalance := broker.GetCapital()
	accountRisk := accountBalance * (oc.riskPerTradePercent / 100)

	entryPrice := history.GetPrice()
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
}

func (oc *capitalAllocatorFixedRisk) Format() *formatter.FormatterNode {
	return formatter.Format("CapitalAllocatorFixedRisk",
		formatter.Format(fmt.Sprintf("RiskPerTradePercent: %.2f%%", oc.riskPerTradePercent)),
	)
}
