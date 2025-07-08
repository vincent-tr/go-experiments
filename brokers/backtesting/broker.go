package backtesting

import (
	"fmt"
	"go-experiments/brokers"
	"go-experiments/common"
	"time"
)

var log = common.NewLogger("backtesting")

type broker struct {
	ticks         []Tick
	currentIndex  int
	capital       float64
	openPositions []*brokers.Position
	callbacks     map[brokers.Timeframe][]func(candle brokers.Candle)
}

// Run implements brokers.BacktestingBroker.
func (b *broker) Run() error {
	panic("unimplemented")
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
	var price float64
	switch order.Direction {
	case brokers.PositionDirectionLong:
		price = b.currentTick().Ask
	case brokers.PositionDirectionShort:
		price = b.currentTick().Bid
	default:
		return nil, fmt.Errorf("invalid position direction: %s", order.Direction)
	}

	panic("unimplemented")
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
	log.Debug("📊 Read %d ticks from CSV file.", len(ticks))

	b := &broker{
		ticks:         ticks,
		currentIndex:  0,
		capital:       initialCapital,
		openPositions: make([]*brokers.Position, 0),
		callbacks:     make(map[brokers.Timeframe][]func(candle brokers.Candle)),
	}

	return b, nil
}

func (b *broker) currentTick() Tick {
	return b.ticks[b.currentIndex]
}
