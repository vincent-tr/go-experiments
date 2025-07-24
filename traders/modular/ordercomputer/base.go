package ordercomputer

import (
	"go-experiments/brokers"
	"go-experiments/traders/modular/context"
	"go-experiments/traders/modular/formatter"
)

type OrderComputer interface {
	formatter.Formatter
	Compute(ctx context.TraderContext, order *brokers.Order) error
}

func newOrderComputer(
	compute func(ctx context.TraderContext, order *brokers.Order) error,
	format func() *formatter.FormatterNode,
) OrderComputer {
	return &orderComputer{
		compute: compute,
		format:  format,
	}
}

type orderComputer struct {
	compute func(ctx context.TraderContext, order *brokers.Order) error
	format  func() *formatter.FormatterNode
}

func (oc *orderComputer) Compute(ctx context.TraderContext, order *brokers.Order) error {
	return oc.compute(ctx, order)
}
func (oc *orderComputer) Format() *formatter.FormatterNode {
	return oc.format()
}
