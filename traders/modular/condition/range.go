package condition

import (
	"fmt"
	"go-experiments/traders/modular/context"
	"go-experiments/traders/modular/formatter"
	"go-experiments/traders/modular/indicators"
)

func IndicatorRange(indicator indicators.Indicator, min, max float64) Condition {
	return newCondition(
		func(ctx context.TraderContext) bool {
			values := indicator.Values(ctx)
			if len(values) == 0 {
				return false
			}
			lastValue := values[len(values)-1]
			return lastValue >= min && lastValue <= max
		},
		func() *formatter.FormatterNode {
			return formatter.Format("IndicatorRange",
				indicator.Format(),
				formatter.Format(fmt.Sprintf("Min: %.2f", min)),
				formatter.Format(fmt.Sprintf("Max: %.2f", max)),
			)
		},
	)
}
