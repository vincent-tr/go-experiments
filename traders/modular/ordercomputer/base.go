package ordercomputer

import (
	"go-experiments/brokers"
	"go-experiments/traders/modular/formatter"
	"go-experiments/traders/tools"
)

type OrderComputer interface {
	formatter.Formatter
	Compute(broker brokers.Broker, history *tools.History, order *brokers.Order) error
}
