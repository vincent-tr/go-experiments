package modular

import (
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
