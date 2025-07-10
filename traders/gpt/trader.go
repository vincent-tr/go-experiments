package gpt

import (
	"go-experiments/brokers"
	"go-experiments/common"
	"go-experiments/traders/tools"

	"github.com/markcheno/go-talib"
)

var log = common.NewLogger("traders/gpt")

func Setup(broker brokers.Broker) {

	trader := newTrader(broker)

	broker.RegisterMarketDataCallback(brokers.Timeframe1Minute, func(candle brokers.Candle) {
		trader.tick(candle)
	})
}

type trader struct {
	broker  brokers.Broker
	history *tools.History
}

func newTrader(broker brokers.Broker) *trader {
	return &trader{
		broker:  broker,
		history: tools.NewHistory(100), // at least 21 to have prev EMA20. 100 to have enough data for last waves
	}
}

func (t *trader) tick(candle brokers.Candle) {
	t.history.AddCandle(candle)

	res, direction := t.shouldTakePosition()
	if !res {
		return
	}

}

func (t *trader) shouldTakePosition() (bool, brokers.PositionDirection) {
	var defaultValue brokers.PositionDirection

	closePrices := t.history.GetClosePrices()
	if closePrices == nil {
		log.Debug("Not enough data")
		return false, defaultValue
	}

	ema20 := talib.Ema(closePrices, 20)
	ema5 := talib.Ema(closePrices, 5)
	rsi := talib.Rsi(closePrices, 14)
	last := len(closePrices) - 1

	prevFast := ema5[last-1]
	prevSlow := ema20[last-1]
	currFast := ema5[last]
	currSlow := ema20[last]
	currRSI := rsi[last]

	// RSI must be between 30 and 70 (neutral zone)
	if currRSI < 30 || currRSI > 70 {
		return false, defaultValue
	}

	// Buy signal: bullish crossover
	if prevFast < prevSlow && currFast > currSlow {
		return true, brokers.PositionDirectionLong
	}

	// Sell signal: bearish crossover
	if prevFast > prevSlow && currFast < currSlow {
		return true, brokers.PositionDirectionShort
	}

	return false, defaultValue
}
