package backtesting

import (
	"go-experiments/brokers"
	"time"
)

type position struct {
	// Open position details
	direction brokers.PositionDirection
	quantity  int
	openPrice float64
	openTime  time.Time
	capital   float64 // Account capital at the time of opening

	// Close trigger details
	stopLoss   float64
	takeProfit float64

	// Close position details
	closePrice float64
	closeTime  time.Time
	closed     bool
}

// Direction implements brokers.Position.
func (p *position) Direction() brokers.PositionDirection {
	return p.direction
}

// Quantity implements brokers.Position.
func (p *position) Quantity() int {
	return p.quantity
}

// OpenPrice implements brokers.Position.
func (p *position) OpenPrice() float64 {
	return p.openPrice
}

// OpenTime implements brokers.Position.
func (p *position) OpenTime() time.Time {
	return p.openTime
}

// ClosePrice implements brokers.Position.
func (p *position) ClosePrice() float64 {
	return p.closePrice
}

// CloseTime implements brokers.Position.
func (p *position) CloseTime() time.Time {
	return p.closeTime
}

// Closed implements brokers.Position.
func (p *position) Closed() bool {
	return p.closed
}

var _ brokers.Position = (*position)(nil)

func newPosition(currentTick *tick, capital float64, order *brokers.Order) *position {

	return &position{
		direction: order.Direction,
		quantity:  order.Quantity,
		openPrice: getOpenPrice(order.Direction, currentTick),
		openTime:  currentTick.Timestamp,
		capital:   capital,

		stopLoss:   order.StopLoss,
		takeProfit: order.TakeProfit,
	}
}

type CloseTrigger int

const (
	CloseTriggerNone CloseTrigger = iota
	CloseTriggerStopLoss
	CloseTriggerTakeProfit
)

// isTriggered checks if the position should be closed based on the current tick.
func (pos *position) isTriggered(currentTick *tick) CloseTrigger {
	price := getClosePrice(pos.direction, currentTick)

	switch pos.direction {

	case brokers.PositionDirectionLong:
		// For long positions, we check if the price is below the stop loss or above the take profit.
		if price <= pos.stopLoss {
			return CloseTriggerStopLoss
		}
		if price >= pos.takeProfit {
			return CloseTriggerTakeProfit
		}

		return CloseTriggerNone

	case brokers.PositionDirectionShort:
		// For short positions, we check if the price is above the stop loss or below the take profit.
		if price >= pos.stopLoss {
			return CloseTriggerStopLoss
		}
		if price <= pos.takeProfit {
			return CloseTriggerTakeProfit
		}

		return CloseTriggerNone

	default:
		panic("invalid position direction: " + pos.direction.String())
	}
}

func (pos *position) closePosition(currentTick *tick) {
	pos.closePrice = getClosePrice(pos.direction, currentTick)
	pos.closeTime = currentTick.Timestamp
	pos.closed = true
}

func getOpenPrice(direction brokers.PositionDirection, currentTick *tick) float64 {
	switch direction {

	case brokers.PositionDirectionLong:
		return currentTick.Ask

	case brokers.PositionDirectionShort:
		return currentTick.Bid

	default:
		panic("invalid position direction: " + direction.String())
	}
}

func getClosePrice(direction brokers.PositionDirection, currentTick *tick) float64 {
	switch direction {

	case brokers.PositionDirectionLong:
		// For long positions, we use the ask price to close
		// because we buy at the ask price and sell at the bid price.
		// This is the price at which we can close the position.
		return currentTick.Bid

	case brokers.PositionDirectionShort:
		// For short positions, we use the bid price to close
		// because we sell at the bid price and buy at the ask price.
		// This is the price at which we can close the position.
		// Note that this is the opposite of the open price for long positions.
		return currentTick.Ask

	default:
		panic("invalid position direction: " + direction.String())
	}
}

func (pos *position) getMargin(leverage float64) float64 {
	totalAmount := float64(pos.Quantity()) * pos.openPrice
	margin := totalAmount / leverage

	return margin
}

func (pos *position) getProfitOrLoss() float64 {
	if !pos.closed {
		return 0.0
	}

	diff := pos.closePrice - pos.openPrice
	if pos.direction == brokers.PositionDirectionShort {
		diff = -diff
	}
	totalAmount := float64(pos.Quantity()) * diff
	return totalAmount
}
