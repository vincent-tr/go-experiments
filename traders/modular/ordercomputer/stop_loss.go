package ordercomputer

import (
	"fmt"
	"go-experiments/brokers"
	"go-experiments/traders/modular/context"
	"go-experiments/traders/modular/formatter"
	"go-experiments/traders/tools"

	"github.com/markcheno/go-talib"
)

func StopLossAtr(period int, multiplier float64) OrderComputer {
	return &stopLossAtr{
		period:     period,
		multiplier: multiplier,
	}
}

type stopLossAtr struct {
	period     int
	multiplier float64
}

func (oc *stopLossAtr) Compute(broker brokers.Broker, history *tools.History, order brokers.Order) error {
	atr := talib.Atr(history.GetHighPrices(), history.GetLowPrices(), history.GetClosePrices(), oc.period)

	if len(atr) == 0 {
		return fmt.Errorf("not enough data for ATR calculation")
	}

	lastAtr := atr[len(atr)-1]
	order.StopLoss = lastAtr * oc.multiplier
	return nil
}

func (oc *stopLossAtr) Format() *formatter.FormatterNode {
	return formatter.Format("StopLossATR",
		formatter.Format(fmt.Sprintf("Period: %d", oc.period)),
		formatter.Format(fmt.Sprintf("Multiplier: %.4f", oc.multiplier)),
	)
}

func StopLossPipBuffer(pipBuffer int, lookupPeriod int) OrderComputer {
	return &stopLossPipBuffer{
		pipBuffer:    pipBuffer,
		lookupPeriod: lookupPeriod,
	}
}

type stopLossPipBuffer struct {
	pipBuffer    int
	lookupPeriod int
}

func (oc *stopLossPipBuffer) Compute(ctx context.TraderContext, order *brokers.Order) error {
	pipSize := 0.0001 // Assuming a pip size of 0.0001 for most currency pairs
	pipDistance := float64(oc.pipBuffer) * pipSize

	switch order.Direction {
	case brokers.PositionDirectionLong:
		// find lowest low in last lookupPeriod minutes
		lowest := ctx.HistoricalData().GetLowest(oc.lookupPeriod)
		order.StopLoss = lowest - pipDistance
		return nil
	case brokers.PositionDirectionShort:
		// find highest high in last lookupPeriod minutes
		highest := ctx.HistoricalData().GetHighest(oc.lookupPeriod)
		order.StopLoss = highest + pipDistance
		return nil
	default:
		return fmt.Errorf("invalid position direction: %s", order.Direction.String())
	}
}

func (oc *stopLossPipBuffer) Format() *formatter.FormatterNode {
	return formatter.Format("StopLossPipBuffer",
		formatter.Format(fmt.Sprintf("Pip Buffer: %d", oc.pipBuffer)),
		formatter.Format(fmt.Sprintf("Lookup Period: %d", oc.lookupPeriod)),
	)
}
