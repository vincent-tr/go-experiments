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
}

// GetCapital implements brokers.Broker.
func (b *broker) GetCapital() float64 {
	return b.capital
}

// GetCurrentTime implements brokers.Broker.
func (b *broker) GetCurrentTime() time.Time {
	panic("unimplemented")
}

// GetMarketDataChannel implements brokers.Broker.
func (b *broker) GetMarketDataChannel(timeframe brokers.Timeframe) <-chan brokers.Candle {
	panic("unimplemented")
}

// PlaceOrder implements brokers.Broker.
func (b *broker) PlaceOrder(order *brokers.Order) (brokers.Position, error) {
	panic("unimplemented")
}

var _ brokers.Broker = (*broker)(nil)

// NewBroker creates a new instance of the broker.
func NewBroker(beginDate, endDate time.Time, symbol string, initialCapital float64) (brokers.Broker, error) {
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
	}

	return b, nil
}
