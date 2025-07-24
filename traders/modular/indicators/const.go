package indicators

import (
	"fmt"
	"go-experiments/traders/modular/context"
	"go-experiments/traders/modular/formatter"
)

func Const(period int, value float64) Indicator {
	return newIndicator(
		func(ctx context.TraderContext) []float64 {
			values := make([]float64, period)

			for i := range values {
				values[i] = value
			}

			return values
		},
		func() *formatter.FormatterNode {
			return formatter.Format("Const",
				formatter.Format(fmt.Sprintf("Period: %d", period)),
				formatter.Format(fmt.Sprintf("Value: %.4f", value)),
			)
		},
	)
}
