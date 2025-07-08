package backtesting

import (
	"fmt"
	"go-experiments/brokers"
	"go-experiments/common"
	"slices"
	"time"
)

var log = common.NewLogger("backtesting")

type broker struct {
	ticks         []tick
	currentIndex  int
	capital       float64
	openPositions []*position
	callbacks     map[brokers.Timeframe][]func(candle brokers.Candle)
}

// Run implements brokers.BacktestingBroker.
func (b *broker) Run() error {
	log.Info("🚀 Starting backtest with %d ticks and initial capital %.2f", len(b.ticks), b.capital)

	for b.currentIndex < len(b.ticks) {
		b.processTick()
		b.currentIndex++
	}

	common.ClearCurrentTime()

	log.Info("✅ Backtest completed. Final capital: %.2f", b.capital)

	return nil
}

// GetLotSize implements brokers.Broker.
func (b *broker) GetLotSize() int {
	// For backtesting, we assume a lot size of 1 for simplicity.
	// In a real broker, this would be the number of units per lot.
	// Not that using IG broker, EUR/USD Mini has also a size of 1.
	return 1
}

// GetCapital implements brokers.Broker.
func (b *broker) GetCapital() float64 {
	return b.capital
}

// GetCurrentTime implements brokers.Broker.
func (b *broker) GetCurrentTime() time.Time {
	return b.currentTick().Timestamp
}

// GetMarketDataChannel implements brokers.Broker.
func (b *broker) RegisterMarketDataCallback(timeframe brokers.Timeframe, callback func(candle brokers.Candle)) {
	b.callbacks[timeframe] = append(b.callbacks[timeframe], callback)
}

// PlaceOrder implements brokers.Broker.
func (b *broker) PlaceOrder(order *brokers.Order) (brokers.Position, error) {
	position := newPosition(b.currentTick(), order)

	// Calculate the total amount of money invested based on the lot size and quantity.
	totalAmount := float64(order.Quantity*b.GetLotSize()) * position.openPrice
	if totalAmount > b.capital {
		return nil, fmt.Errorf("insufficient capital: cannot place order for %d lots at price %.2f (total: %.2f)", order.Quantity, position.openPrice, totalAmount)
	}

	b.capital -= totalAmount
	b.openPositions = append(b.openPositions, position)

	return position, nil
}

var _ brokers.Broker = (*broker)(nil)
var _ brokers.BacktestingBroker = (*broker)(nil)

// NewBroker creates a new instance of the broker.
func NewBroker(beginDate, endDate time.Time, symbol string, initialCapital float64) (brokers.BacktestingBroker, error) {
	beginTime := time.Now()

	ticks, err := loadData(beginDate, endDate, symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to load data: %w", err)
	}

	endTime := time.Now()
	duration := endTime.Sub(beginTime)
	log.Debug("⏱️ Unzipped and parsed CSV in %s.", duration)
	log.Debug("📊 Read %d ticks from CSV.", len(ticks))

	b := &broker{
		ticks:         ticks,
		currentIndex:  0,
		capital:       initialCapital,
		openPositions: make([]*position, 0),
		callbacks:     make(map[brokers.Timeframe][]func(candle brokers.Candle)),
	}

	return b, nil
}

func (b *broker) currentTick() *tick {
	return &b.ticks[b.currentIndex]
}

func (b *broker) processTick() {
	currentTick := b.currentTick()
	common.SetCurrentTime(currentTick.Timestamp)

	log.Debug("📈 Processing tick at %s: Bid=%.5f, Ask=%.5f", currentTick.Timestamp.Format("2006-01-02 15:04:05"), currentTick.Bid, currentTick.Ask)

	for _, pos := range b.openPositions {
		if pos.isTriggered(currentTick) {
			pos.closePosition(currentTick)
			log.Debug("📉 Position closed at %s: Direction=%s, Quantity=%d, OpenPrice=%.5f, ClosePrice=%.5f",
				currentTick.Timestamp.Format("2006-01-02 15:04:05"), pos.direction, pos.quantity, pos.openPrice, pos.closePrice)
		}
	}

	// Check if we have a full candle for any registered timeframes
	for timeframe, callbacks := range b.callbacks {
		candle := b.tryCandle(timeframe)

		if candle != nil {
			log.Debug("📊 New candle for timeframe %s: Open=%.5f, Close=%.5f, High=%.5f, Low=%.5f",
				timeframe, candle.Open, candle.Close, candle.High, candle.Low)

			// Call all registered callbacks for this timeframe
			for _, callback := range callbacks {
				callback(*candle)
			}
		}
	}

}

func (b *broker) tryCandle(timeframe brokers.Timeframe) *brokers.Candle {
	currentTick := b.currentTick()

	var nextTick *tick
	if b.currentIndex+1 < len(b.ticks) {
		nextTick = &b.ticks[b.currentIndex+1]
	}

	currentBucket := getTimeframeBucket(currentTick, timeframe)

	// Is the next tick in a different timeframe?
	if nextTick != nil && getTimeframeBucket(nextTick, timeframe) == currentBucket {
		// No complete candle yet, we need to wait for the next tick
		return nil
	}

	// We have the last tick of the current timeframe

	// Get all ticks for the current timeframe bucket
	timeframeTicks := []*tick{currentTick}
	for i := b.currentIndex - 1; i >= 0; i-- {
		if getTimeframeBucket(&b.ticks[i], timeframe) == currentBucket {
			timeframeTicks = append(timeframeTicks, &b.ticks[i])
		} else {
			break
		}
	}

	slices.Reverse(timeframeTicks)

	// Create a candle from the timeframe ticks
	low := timeframeTicks[0].Price()
	high := timeframeTicks[0].Price()
	for _, t := range timeframeTicks {
		price := t.Price()
		if price < low {
			low = price
		}
		if price > high {
			high = price
		}
	}

	return &brokers.Candle{
		Open:  timeframeTicks[0].Price(),
		Close: timeframeTicks[len(timeframeTicks)-1].Price(),
		High:  high,
		Low:   low,
	}
}

func getTimeframeBucket(tick *tick, timeframe brokers.Timeframe) string {
	// This function should return the start time of the bucket for the given timeframe.
	// For simplicity, we assume that the tick's timestamp is already aligned with the timeframe.
	return tick.Timestamp.Truncate(time.Duration(timeframe)).Format("2006-01-02 15:04:05")
}
