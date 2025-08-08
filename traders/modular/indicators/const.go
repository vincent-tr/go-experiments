package indicators

import (
	"encoding/json"
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

func init() {
	jsonParsers.RegisterParser("const", func(arg json.RawMessage) (Indicator, error) {
		var params struct {
			Period int     `json:"period"`
			Value  float64 `json:"value"`
		}

		if err := json.Unmarshal(arg, &params); err != nil {
			return nil, fmt.Errorf("failed to parse Const parameters: %w", err)
		}

		return Const(params.Period, params.Value), nil
	})
}
