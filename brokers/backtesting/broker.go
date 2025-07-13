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
	lotSize          int
	leverage         float64
	ticks            []tick
	currentIndex     int
	capital          float64
	openPositions    map[*position]struct{}
	callbacks        map[brokers.Timeframe][]func(candle brokers.Candle)
	positionsHistory []*position
}

// Run implements brokers.BacktestingBroker.
func (b *broker) Run() error {
	log.Info("🚀 Starting backtest with %d ticks and initial capital %.2f", len(b.ticks), b.capital)

	for {
		b.processTick()

		if b.currentIndex == len(b.ticks)-1 {
			break
		}

		b.currentIndex++
	}

	b.closeAllOpenPositions()

	common.ClearCurrentTime()

	log.Info("✅ Backtest completed.")
	b.printSummary()

	return nil
}

// GetLotSize implements brokers.Broker.
func (b *broker) GetLotSize() int {
	return b.lotSize
}

// GetCapital implements brokers.Broker.
func (b *broker) GetCapital() float64 {
	return b.capital
}

// GetLeverage implements brokers.Broker.
func (b *broker) GetLeverage() float64 {
	return b.leverage
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
	pos := newPosition(b.currentTick(), b.GetCapital(), order)
	margin := pos.getMargin(b.GetLeverage())

	if margin > b.capital {
		return nil, fmt.Errorf("insufficient capital: cannot place order for %d lots at price %.4f (margin: %.2f, capital:  %.2f)", pos.Quantity(), pos.OpenPrice(), margin, b.capital)
	}

	b.capital -= margin
	b.openPositions[pos] = struct{}{}
	b.positionsHistory = append(b.positionsHistory, pos)

	log.Debug("📈 Placed order: Direction=%s, Quantity=%d, OpenPrice=%.5f, StopLoss=%.5f, TakeProfit=%.5f, Reason=%s",
		pos.Direction(), pos.Quantity(), pos.openPrice, order.StopLoss, order.TakeProfit,
		order.Reason)

	return pos, nil
}

var _ brokers.Broker = (*broker)(nil)
var _ brokers.BacktestingBroker = (*broker)(nil)

// NewBroker creates a new instance of the broker.
func NewBroker(beginDate, endDate time.Time, symbol string, lotSize int, leverage float64, initialCapital float64) (brokers.BacktestingBroker, error) {
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
		lotSize:          lotSize,
		leverage:         leverage,
		ticks:            ticks,
		currentIndex:     0,
		capital:          initialCapital,
		openPositions:    make(map[*position]struct{}),
		callbacks:        make(map[brokers.Timeframe][]func(candle brokers.Candle)),
		positionsHistory: make([]*position, 0),
	}

	return b, nil
}

func (b *broker) currentTick() *tick {
	return &b.ticks[b.currentIndex]
}

func (b *broker) processTick() {
	currentTick := b.currentTick()
	common.SetCurrentTime(currentTick.Timestamp)

	// log.Debug("📈 Processing tick at %s: Bid=%.5f, Ask=%.5f", currentTick.Timestamp.Format("2006-01-02 15:04:05"), currentTick.Bid, currentTick.Ask)

	for pos := range b.openPositions {
		switch pos.isTriggered(currentTick) {
		case CloseTriggerNone:
			// Position is still open, do nothing
			continue
		case CloseTriggerStopLoss, CloseTriggerTakeProfit:
			// Position should be closed
			b.closePosition(pos)

			closeReason := "unknown"
			if pos.isTriggered(currentTick) == CloseTriggerStopLoss {
				closeReason = "stop loss"
			} else if pos.isTriggered(currentTick) == CloseTriggerTakeProfit {
				closeReason = "take profit"
			}

			log.Debug("📉 Position closed (%s) at %s: Direction=%s, Quantity=%d, OpenPrice=%.5f, ClosePrice=%.5f",
				closeReason,
				currentTick.Timestamp.Format("2006-01-02 15:04:05"),
				pos.direction, pos.quantity, pos.openPrice, pos.closePrice)
		}
	}

	// Check if we have a full candle for any registered timeframes
	for timeframe, callbacks := range b.callbacks {
		candle := b.tryCandle(timeframe)

		if candle != nil {
			// log.Debug("📊 New candle for timeframe %s: Open=%.5f, Close=%.5f, High=%.5f, Low=%.5f",
			// 	timeframe, candle.Open, candle.Close, candle.High, candle.Low)

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

func (b *broker) closeAllOpenPositions() {
	for pos := range b.openPositions {
		b.closePosition(pos)

		log.Debug("📉 Position closed (end of test) at %s: Direction=%s, Quantity=%d, OpenPrice=%.5f, ClosePrice=%.5f",
			b.currentTick().Timestamp.Format("2006-01-02 15:04:05"),
			pos.direction, pos.quantity, pos.openPrice, pos.closePrice)
	}
}

func (b *broker) closePosition(pos *position) {
	pos.closePosition(b.currentTick())
	delete(b.openPositions, pos)

	b.capital += pos.getMargin(b.GetLeverage())
	b.capital += pos.getProfitOrLoss()
}

func (b *broker) printSummary() {
	log.Info("📊 Backtest Summary:")
	log.Info("Total positions: %d", len(b.positionsHistory))
	log.Info("Final capital: %.2f", b.capital)

	log.Info("Positions history:")
	for _, pos := range b.positionsHistory {
		profit := pos.getProfitOrLoss()
		var profitColor string
		if profit < 0 {
			profitColor = fmt.Sprintf("\033[31m%.2f\033[0m", profit) // Red for losses
		} else if profit > 0 {
			profitColor = fmt.Sprintf("\033[32m%.2f\033[0m", profit) // Green for profits
		} else {
			profitColor = fmt.Sprintf("%.2f", profit) // No color for zero
		}

		log.Info(" - Capital: %0.2f, Direction: %s, OpenTime: %s, Profit: %s, Duration: %s",
			pos.capital, pos.direction, pos.openTime.Format("2006-01-02 15:04:05"), profitColor, pos.CloseTime().Sub(pos.OpenTime()).String())
	}

	log.Info("Total capital at the end of the test: %.2f", b.capital)
	log.Info("Total number of positions taken: %d", len(b.positionsHistory))
}
