package modular

import (
	"fmt"
	"go-experiments/brokers"
	"go-experiments/common"
	"go-experiments/traders/modular/condition"
	"go-experiments/traders/modular/context"
	"go-experiments/traders/modular/ordercomputer"
	"go-experiments/traders/tools"
	"maps"
	"slices"
	"time"
)

var log = common.NewLogger("traders/modular")

func Setup(broker brokers.Broker, builder Builder) error {

	trader, err := newTrader(broker, builder)
	if err != nil {
		return err
	}

	broker.RegisterMarketDataCallback(brokers.Timeframe1Minute, func(candle brokers.Candle) {
		trader.tick(candle)
	})

	return nil
}

type trader struct {
	broker           brokers.Broker
	history          *tools.History
	openPositions    map[brokers.Position]struct{}
	filter           condition.Condition
	longTrigger      condition.Condition
	shortTrigger     condition.Condition
	stopLoss         ordercomputer.OrderComputer
	takeProfit       ordercomputer.OrderComputer
	capitalAllocator ordercomputer.OrderComputer
}

func newTrader(broker brokers.Broker, builder Builder) (*trader, error) {
	b, err := getBuilder(builder)
	if err != nil {
		return nil, err
	}

	if b.historySize <= 0 {
		return nil, fmt.Errorf("history size must be greater than 0")
	}
	if b.filter == nil {
		return nil, fmt.Errorf("filter must be set")
	}
	if b.longTrigger == nil && b.shortTrigger == nil {
		return nil, fmt.Errorf("either long or short trigger must be set")
	}
	if b.stopLoss == nil {
		return nil, fmt.Errorf("stop loss computer must be set")
	}
	if b.takeProfit == nil {
		return nil, fmt.Errorf("take profit computer must be set")
	}
	if b.capitalAllocator == nil {
		return nil, fmt.Errorf("capital allocator must be set")
	}

	return &trader{
		broker:           broker,
		history:          tools.NewHistory(b.historySize),
		openPositions:    make(map[brokers.Position]struct{}),
		filter:           b.filter,
		longTrigger:      b.longTrigger,
		shortTrigger:     b.shortTrigger,
		stopLoss:         b.stopLoss,
		takeProfit:       b.takeProfit,
		capitalAllocator: b.capitalAllocator,
	}, nil
}

// workaround builder naming conflict
func getBuilder(bi Builder) (*builder, error) {
	b, ok := bi.(*builder)
	if !ok {
		return nil, fmt.Errorf("invalid builder type: %T", bi)
	}
	return b, nil
}

func (t *trader) tick(candle brokers.Candle) {
	t.history.AddCandle(candle)

	for pos := range t.openPositions {
		if pos.Closed() {
			delete(t.openPositions, pos)
		}
	}

	if !t.filter.Execute(t) {
		return
	}

	shouldTakeLong := t.longTrigger.Execute(t)
	shouldTakeShort := t.shortTrigger.Execute(t)

	if shouldTakeLong && shouldTakeShort {
		log.Warning("Both long and short triggers are true, ignoring")
		return
	}

	if shouldTakeLong {
		t.takePosition(brokers.PositionDirectionLong)
	}

	if shouldTakeShort {
		t.takePosition(brokers.PositionDirectionShort)
	}
}

func (t *trader) takePosition(direction brokers.PositionDirection) {
	order := &brokers.Order{
		Direction: direction,
	}

	err := t.stopLoss.Compute(t, order)
	if err != nil {
		log.Error("Failed to compute stop loss: %v", err)
		return
	}

	err = t.takeProfit.Compute(t, order)
	if err != nil {
		log.Error("Failed to compute take profit: %v", err)
		return
	}

	err = t.capitalAllocator.Compute(t, order)
	if err != nil {
		log.Error("Failed to compute capital allocation: %v", err)
		return
	}

	// TODO: reason

	pos, err := t.broker.PlaceOrder(order)
	if err != nil {
		log.Error("Failed to place order: %v", err)
		return
	}

	t.openPositions[pos] = struct{}{}
}

var _ context.TraderContext = (*trader)(nil)

func (t *trader) Broker() brokers.Broker {
	return t.broker
}
func (t *trader) HistoricalData() *tools.History {
	return t.history
}

func (t *trader) OpenPositions() []brokers.Position {
	return slices.Collect(maps.Keys(t.openPositions))
}

func (t *trader) Timestamp() time.Time {
	return t.broker.GetCurrentTime()
}

func (t *trader) EntryPrice() float64 {
	return t.history.GetPrice()
}
