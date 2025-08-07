package ordercomputer

import (
	"go-experiments/brokers"
	"go-experiments/traders/modular/context"
	"go-experiments/traders/modular/formatter"
	"go-experiments/traders/modular/json"
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

var jsonParsers = json.NewRegistry[OrderComputer]()

func FromJSON(jsonData []byte) (OrderComputer, error) {
	return jsonParsers.FromJSON(jsonData)
}

func ToJSON(oc OrderComputer) ([]byte, error) {
	panic("ToJSON not implemented for order computers")
}
