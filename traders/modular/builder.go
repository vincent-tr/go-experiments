package modular

import (
	"go-experiments/traders/modular/formatter"
	"go-experiments/traders/modular/marketcondition"
	"go-experiments/traders/modular/opentimecondition"
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
	SetOpenTimeCondition(condition opentimecondition.OpenTimeCondition) StrategyBuilder
	SetFilter(filter marketcondition.MarketCondition) StrategyBuilder
	SetLongTrigger(trigger marketcondition.MarketCondition) StrategyBuilder
	SetShortTrigger(trigger marketcondition.MarketCondition) StrategyBuilder
}

type RiskManagerBuilder interface {
	SetStopLoss(computer ordercomputer.OrderComputer) RiskManagerBuilder
	SetTakeProfit(computer ordercomputer.OrderComputer) RiskManagerBuilder
}

type CapitalAllocatorBuilder interface {
	SetAllocator(computer ordercomputer.OrderComputer) CapitalAllocatorBuilder
}

type builder struct {
	historySize       int
	openTimeCondition opentimecondition.OpenTimeCondition
	filter            marketcondition.MarketCondition
	longTrigger       marketcondition.MarketCondition
	shortTrigger      marketcondition.MarketCondition
	stopLoss          ordercomputer.OrderComputer
	takeProfit        ordercomputer.OrderComputer
	capitalAllocator  ordercomputer.OrderComputer
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

func (b *builder) SetOpenTimeCondition(condition opentimecondition.OpenTimeCondition) StrategyBuilder {
	b.openTimeCondition = condition
	return b
}

func (b *builder) SetFilter(filter marketcondition.MarketCondition) StrategyBuilder {
	b.filter = filter
	return b
}

func (b *builder) SetLongTrigger(trigger marketcondition.MarketCondition) StrategyBuilder {
	b.longTrigger = trigger
	return b
}

func (b *builder) SetShortTrigger(trigger marketcondition.MarketCondition) StrategyBuilder {
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
