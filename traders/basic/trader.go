package basic

import (
	"go-experiments/brokers"
	"go-experiments/common"
)

var log = common.NewLogger("traders/basic")

func Setup(broker brokers.Broker) {
	broker.RegisterMarketDataCallback(brokers.Timeframe1Minute, func(candle brokers.Candle) {

		// Example logic: if the candle closed higher than it opened, place a long order
		if candle.Close > candle.Open {
			diff := candle.Close - candle.Open

			order := &brokers.Order{
				Direction:  brokers.PositionDirectionLong,
				Quantity:   10,
				StopLoss:   candle.Low,
				TakeProfit: candle.Close + diff*2,
				Reason:     "BasicTrader: Long on bullish candle",
			}

			if _, err := broker.PlaceOrder(order); err != nil {
				log.Error("Failed to place order: %v", err)
			}
		}
	})
}
