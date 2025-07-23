package modular

import (
	"go-experiments/traders/modular/condition"
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
	SetFilter(condition condition.Condition) StrategyBuilder
	SetLongTrigger(trigger condition.Condition) StrategyBuilder
	SetShortTrigger(trigger condition.Condition) StrategyBuilder
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
	filter           condition.Condition
	longTrigger      condition.Condition
	shortTrigger     condition.Condition
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

func (b *builder) SetFilter(filter condition.Condition) StrategyBuilder {
	b.filter = filter
	return b
}

func (b *builder) SetLongTrigger(trigger condition.Condition) StrategyBuilder {
	b.longTrigger = trigger
	return b
}

func (b *builder) SetShortTrigger(trigger condition.Condition) StrategyBuilder {
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
	panic("TODO: implement builder.Format()")
}
