package modular

import (
	"fmt"
	"go-experiments/brokers"
	"go-experiments/common"
	"go-experiments/traders/modular/marketcondition"
	"go-experiments/traders/modular/opentimecondition"
	"go-experiments/traders/modular/ordercomputer"
	"go-experiments/traders/tools"
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
	broker            brokers.Broker
	history           *tools.History
	openPosition      map[brokers.Position]struct{}
	openTimeCondition opentimecondition.OpenTimeCondition
	filter            marketcondition.MarketCondition
	longTrigger       marketcondition.MarketCondition
	shortTrigger      marketcondition.MarketCondition
	stopLoss          ordercomputer.OrderComputer
	takeProfit        ordercomputer.OrderComputer
	capitalAllocator  ordercomputer.OrderComputer
}

func newTrader(broker brokers.Broker, builder Builder) (*trader, error) {
	b, err := getBuilder(builder)
	if err != nil {
		return nil, err
	}

	if b.historySize <= 0 {
		return nil, fmt.Errorf("history size must be greater than 0")
	}
	if b.openTimeCondition == nil {
		return nil, fmt.Errorf("open time condition must be set")
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
		broker:            broker,
		history:           tools.NewHistory(b.historySize),
		openPosition:      make(map[brokers.Position]struct{}),
		openTimeCondition: b.openTimeCondition,
		filter:            b.filter,
		longTrigger:       b.longTrigger,
		shortTrigger:      b.shortTrigger,
		stopLoss:          b.stopLoss,
		takeProfit:        b.takeProfit,
		capitalAllocator:  b.capitalAllocator,
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

	for pos := range t.openPosition {
		if pos.Closed() {
			delete(t.openPosition, pos)
		}
	}

	// TODO:
	// // Only take one position at a time
	// if t.openPosition != nil {
	// 	return
	// }

	currentTime := t.broker.GetCurrentTime()

	if !t.openTimeCondition.Execute(currentTime) {
		return
	}

	if !t.filter.Execute(t.history) {
		return
	}

	shouldTakeLong := t.longTrigger.Execute(t.history)
	shouldTakeShort := t.shortTrigger.Execute(t.history)

	if shouldTakeLong && shouldTakeShort {
		log.Warning("Both long and short triggers are true, ignoring")
		return
	}

	if shouldTakeLong {
		t.takePosition(candle, brokers.PositionDirectionLong)
	}

	if shouldTakeShort {
		t.takePosition(candle, brokers.PositionDirectionShort)
	}
}

func (t *trader) takePosition(candle brokers.Candle, direction brokers.PositionDirection) {
	order := &brokers.Order{
		Direction: direction,
	}

	err := t.stopLoss.Compute(t.broker, t.history, order)
	if err != nil {
		log.Error("Failed to compute stop loss: %v", err)
		return
	}

	err = t.takeProfit.Compute(t.broker, t.history, order)
	if err != nil {
		log.Error("Failed to compute take profit: %v", err)
		return
	}

	err = t.capitalAllocator.Compute(t.broker, t.history, order)
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

	t.openPosition[pos] = struct{}{}
}
