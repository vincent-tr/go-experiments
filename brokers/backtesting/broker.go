package backtesting

import (
	"fmt"
	"go-experiments/brokers"
	"go-experiments/common"
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
	panic("unimplemented")
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
