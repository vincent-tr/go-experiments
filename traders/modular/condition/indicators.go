package condition

import (
	"fmt"
	"go-experiments/traders/modular/context"
	"go-experiments/traders/modular/formatter"

	"github.com/markcheno/go-talib"
)

func RsiRange(period int, min, max float64) Condition {
	return &rsiRangeCondition{
		period: period,
		min:    min,
		max:    max,
	}
}

type rsiRangeCondition struct {
	period int
	min    float64
	max    float64
}

func (r *rsiRangeCondition) Execute(ctx context.TraderContext) bool {
	closePrices := ctx.HistoricalData().GetClosePrices()
	rsi := talib.Rsi(closePrices, r.period)

	if len(rsi) == 0 {
		return false
	}

	lastRsi := rsi[len(rsi)-1]
	return lastRsi >= r.min && lastRsi <= r.max
}

func (r *rsiRangeCondition) Format() *formatter.FormatterNode {
	return formatter.Format("RsiRange",
		formatter.Format(fmt.Sprintf("Period: %d", r.period)),
		formatter.Format(fmt.Sprintf("Min: %.2f", r.min)),
		formatter.Format(fmt.Sprintf("Max: %.2f", r.max)),
	)
}
