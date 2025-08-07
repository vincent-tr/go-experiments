package modular

import (
	"encoding/json"
	"fmt"
	"go-experiments/traders/modular/conditions"
	"go-experiments/traders/modular/formatter"
	"go-experiments/traders/modular/ordercomputer"
)

type Builder interface {
	formatter.Formatter
	SetHistorySize(size int) Builder
	Strategy() StrategyBuilder
	RiskManager() RiskManagerBuilder
	CapitalAllocator() CapitalAllocatorBuilder
}

func NewBuilder() Builder {
	return &builder{}
}

type StrategyBuilder interface {
	SetFilter(condition conditions.Condition) StrategyBuilder
	SetLongTrigger(trigger conditions.Condition) StrategyBuilder
	SetShortTrigger(trigger conditions.Condition) StrategyBuilder
}

type RiskManagerBuilder interface {
	SetStopLoss(computer ordercomputer.OrderComputer) RiskManagerBuilder
	SetTakeProfit(computer ordercomputer.OrderComputer) RiskManagerBuilder
}

type CapitalAllocatorBuilder interface {
	SetAllocator(computer ordercomputer.OrderComputer) CapitalAllocatorBuilder
}

type builder struct {
	historySize      int
	filter           conditions.Condition
	longTrigger      conditions.Condition
	shortTrigger     conditions.Condition
	stopLoss         ordercomputer.OrderComputer
	takeProfit       ordercomputer.OrderComputer
	capitalAllocator ordercomputer.OrderComputer
}

var _ Builder = (*builder)(nil)
var _ StrategyBuilder = (*builder)(nil)
var _ RiskManagerBuilder = (*builder)(nil)
var _ CapitalAllocatorBuilder = (*builder)(nil)

func (b *builder) SetHistorySize(size int) Builder {
	b.historySize = size
	return b
}

func (b *builder) Strategy() StrategyBuilder {
	return b
}

func (b *builder) RiskManager() RiskManagerBuilder {
	return b
}

func (b *builder) CapitalAllocator() CapitalAllocatorBuilder {
	return b
}

func (b *builder) SetFilter(filter conditions.Condition) StrategyBuilder {
	b.filter = filter
	return b
}

func (b *builder) SetLongTrigger(trigger conditions.Condition) StrategyBuilder {
	b.longTrigger = trigger
	return b
}

func (b *builder) SetShortTrigger(trigger conditions.Condition) StrategyBuilder {
	b.shortTrigger = trigger
	return b
}

func (b *builder) SetStopLoss(computer ordercomputer.OrderComputer) RiskManagerBuilder {
	b.stopLoss = computer
	return b
}

func (b *builder) SetTakeProfit(computer ordercomputer.OrderComputer) RiskManagerBuilder {
	b.takeProfit = computer
	return b
}

func (b *builder) SetAllocator(computer ordercomputer.OrderComputer) CapitalAllocatorBuilder {
	b.capitalAllocator = computer
	return b
}

func (b *builder) Format() *formatter.FormatterNode {
	return formatter.Format("ModularTrader",
		formatter.Format(fmt.Sprintf("HistorySize: %d", b.historySize)),
		formatter.FormatWithChildren("Filter", b.filter),
		formatter.FormatWithChildren("LongTrigger", b.longTrigger),
		formatter.FormatWithChildren("ShortTrigger", b.shortTrigger),
		formatter.FormatWithChildren("StopLoss", b.stopLoss),
		formatter.FormatWithChildren("TakeProfit", b.takeProfit),
		formatter.FormatWithChildren("CapitalAllocator", b.capitalAllocator),
	)
}

type builderJSON struct {
	HistorySize      int             `json:"historySize"`
	Filter           json.RawMessage `json:"filter"`
	LongTrigger      json.RawMessage `json:"longTrigger"`
	ShortTrigger     json.RawMessage `json:"shortTrigger"`
	StopLoss         json.RawMessage `json:"stopLoss"`
	TakeProfit       json.RawMessage `json:"takeProfit"`
	CapitalAllocator json.RawMessage `json:"capitalAllocator"`
}

func FromJSON(jsonData []byte) (Builder, error) {
	var bjson builderJSON
	err := json.Unmarshal(jsonData, &bjson)
	if err != nil {
		return nil, err
	}

	res := &builder{
		historySize: bjson.HistorySize,
	}

	res.filter, err = conditions.FromJSON(bjson.Filter)
	if err != nil {
		return nil, fmt.Errorf("failed to parse filter condition: %w", err)
	}

	res.longTrigger, err = conditions.FromJSON(bjson.LongTrigger)
	if err != nil {
		return nil, fmt.Errorf("failed to parse long trigger condition: %w", err)
	}

	res.shortTrigger, err = conditions.FromJSON(bjson.ShortTrigger)
	if err != nil {
		return nil, fmt.Errorf("failed to parse short trigger condition: %w", err)
	}

	res.stopLoss, err = ordercomputer.FromJSON(bjson.StopLoss)
	if err != nil {
		return nil, fmt.Errorf("failed to parse stop loss order computer: %w", err)
	}

	res.takeProfit, err = ordercomputer.FromJSON(bjson.TakeProfit)
	if err != nil {
		return nil, fmt.Errorf("failed to parse take profit order computer: %w", err)
	}

	res.capitalAllocator, err = ordercomputer.FromJSON(bjson.CapitalAllocator)
	if err != nil {
		return nil, fmt.Errorf("failed to parse capital allocator order computer: %w", err)
	}

	return res, nil
}

func (b *builder) ToJSON() ([]byte, error) {
	bjson := &builderJSON{
		HistorySize: b.historySize,
	}

	var err error

	bjson.Filter, err = conditions.ToJSON(b.filter)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize filter condition: %w", err)
	}

	bjson.LongTrigger, err = conditions.ToJSON(b.longTrigger)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize long trigger condition: %w", err)
	}

	bjson.ShortTrigger, err = conditions.ToJSON(b.shortTrigger)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize short trigger condition: %w", err)
	}

	bjson.StopLoss, err = ordercomputer.ToJSON(b.stopLoss)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize stop loss order computer: %w", err)
	}

	bjson.TakeProfit, err = ordercomputer.ToJSON(b.takeProfit)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize take profit order computer: %w", err)
	}

	bjson.CapitalAllocator, err = ordercomputer.ToJSON(b.capitalAllocator)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize capital allocator order computer: %w", err)
	}

	return json.Marshal(bjson)
}
