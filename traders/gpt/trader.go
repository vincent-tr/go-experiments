package gpt

import (
	"fmt"
	"go-experiments/brokers"
	"go-experiments/common"
	"go-experiments/traders/tools"
	"math"
	"time"

	"github.com/markcheno/go-talib"
)

var log = common.NewLogger("traders/gpt")

const pipSize = 0.0001

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
		history: tools.NewHistory(30), // at least 21 to have prev EMA20.
	}
}

func (t *trader) tick(candle brokers.Candle) {
	t.history.AddCandle(candle)

	if !t.history.IsComplete() {
		log.Debug("Not enough data to make a decision")
		return
	}

	res, direction := t.shouldTakePosition()
	if !res {
		return
	}

	entryPrice := candle.Close
	stopLoss := t.computeStopLoss(direction)
	takeProfit := t.computeTakeProfit(direction, entryPrice, stopLoss)
	positionSize := t.computePositionSize(stopLoss)

	if positionSize == 0 {
		// Not enough capital to take a position
		return
	}

	order := &brokers.Order{
		Direction:  direction,
		Quantity:   positionSize,
		StopLoss:   stopLoss,
		TakeProfit: takeProfit,
		Reason:     fmt.Sprintf("GPT strategy: %s at %.5f", direction, entryPrice),
	}

	if _, err := t.broker.PlaceOrder(order); err != nil {
		log.Error("Failed to place order: %v", err)
	}
}

func (t *trader) shouldTakePosition() (bool, brokers.PositionDirection) {
	var defaultValue brokers.PositionDirection

	if !t.shouldTrade() {
		return false, defaultValue
	}

	closePrices := t.history.GetClosePrices()

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

func (t *trader) shouldTrade() bool {
	currentTime := t.broker.GetCurrentTime()

	weekday := currentTime.Weekday()
	if weekday < time.Tuesday || weekday > time.Thursday {
		return false
	}

	if common.IsUSHoliday(currentTime) || common.IsUKHoliday(currentTime) {
		return false
	}

	if !common.LondonSession.IsOpen(currentTime) || !common.NYSession.IsOpen(currentTime) {
		return false
	}

	return true
}

// Computes the stop-loss price based on the last 15 minutes of candles.
// For long positions, it is set 3 pips below the lowest low in the last 15 minutes.
// For short positions, it is set 3 pips above the highest high in the last 15 minutes.
func (t *trader) computeStopLoss(direction brokers.PositionDirection) float64 {
	const pipBuffer = 3

	pipDistance := float64(pipBuffer) * pipSize

	switch direction {
	case brokers.PositionDirectionLong:
		// find lowest low in last 15 minutes
		lowest := t.history.GetLowest(15)
		// stop loss is 3 pips below that low
		return lowest - pipDistance

	case brokers.PositionDirectionShort:
		// find highest high in last 15 minutes
		highest := t.history.GetHighest(15)
		// stop loss is 3 pips above that high
		return highest + pipDistance

	default:
		panic("invalid position direction: " + direction.String())
	}
}

// The take-profit is set at a 2:1 reward-to-risk ratio relative to your stop-loss distance.
func (t *trader) computeTakeProfit(direction brokers.PositionDirection, entryPrice, stopLoss float64) float64 {
	switch direction {
	case brokers.PositionDirectionLong:
		risk := entryPrice - stopLoss
		if risk <= 0 {
			panic(fmt.Sprintf("invalid stoploss for long position: entryPrice=%.5f, stopLoss=%.5f", entryPrice, stopLoss))
		}
		return entryPrice + 2*risk

	case brokers.PositionDirectionShort:
		risk := stopLoss - entryPrice
		if risk <= 0 {
			panic(fmt.Sprintf("invalid stoploss for short position: entryPrice=%.5f, stopLoss=%.5f", entryPrice, stopLoss))
		}
		return entryPrice - 2*risk

	default:
		panic("invalid position direction: " + direction.String())
	}
}

func (t *trader) computePositionSize(stopLoss float64) int {
	const riskPercent float64 = 1.0 // 1% risk per trade

	accountBalance := t.broker.GetCapital()
	accountRisk := accountBalance * (riskPercent / 100)

	entryPrice := t.history.GetPrice()
	priceDiff := math.Abs(entryPrice - stopLoss)
	if priceDiff <= 0 {
		panic(fmt.Sprintf("Invalid stop loss price: entryPrice=%.5f, stopLoss=%.5f", entryPrice, stopLoss))
	}

	lotSize := float64(t.broker.GetLotSize())
	riskPerLot := lotSize * priceDiff
	positionSize := accountRisk / riskPerLot

	// Ensure position size doesn't exceed account balance
	// Total value = positionSize * lotSize * entryPrice
	// maxPositionSize := accountBalance / (lotSize * entryPrice)
	// if positionSize > maxPositionSize {
	// 	positionSize = maxPositionSize
	// }

	return int(math.Floor(positionSize))
}
