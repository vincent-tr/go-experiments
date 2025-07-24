package context

import (
	"go-experiments/brokers"
	"go-experiments/traders/tools"
	"time"
)

type TraderContext interface {
	Broker() brokers.Broker
	HistoricalData() *tools.History
	OpenPositions() []brokers.Position
	IndicatorCache() IndicatorCache

	Timestamp() time.Time
	EntryPrice() float64
}

type IndicatorCache interface {
	Tick()
}
