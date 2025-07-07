package backtesting

import (
	"fmt"
	"go-experiments/brokers"
	"log"
	"time"
)

type broker struct {
	ticks        []Tick
	currentIndex int
	capital      float64
}

// GetCapital implements brokers.Broker.
func (b *broker) GetCapital() float64 {
	panic("unimplemented")
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
		log.Fatalf("Failed to load data: %v", err)
	}

	endTime := time.Now()
	duration := endTime.Sub(beginTime)
	fmt.Printf("⏱️ Unzipped and parsed CSV in %s.\n", duration)
	fmt.Printf("📊 Read %d ticks from CSV file.\n", len(ticks))

	fmt.Println("✅ Done.")

	return &broker{}
}
