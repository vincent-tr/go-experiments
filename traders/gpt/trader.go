package gpt

import (
	"go-experiments/brokers"
	"go-experiments/common"
)

var log = common.NewLogger("traders/gpt")

func Setup(broker brokers.Broker) {
	broker.RegisterMarketDataCallback(brokers.Timeframe1Minute, func(candle brokers.Candle) {
		// TODO
	})
}
