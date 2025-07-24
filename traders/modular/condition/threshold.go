package condition

import (
	"fmt"
	"go-experiments/traders/modular/context"
	"go-experiments/traders/modular/formatter"
	"go-experiments/traders/modular/indicators"
)

func Threshold(indicator indicators.Indicator, threshold float64) Condition {
	return newCondition(
		func(ctx context.TraderContext) bool {
			values := indicator.Values(ctx)
			if len(values) == 0 {
				return false
			}
			lastValue := values[len(values)-1]
			return lastValue >= threshold
		},
		func() *formatter.FormatterNode {
			return formatter.Format("Threshold",
				formatter.FormatWithChildren("Indicator", indicator),
				formatter.Format(fmt.Sprintf("Value: %.2f", threshold)),
			)
		},
	)
}
