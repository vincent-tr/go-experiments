package conditions

import (
	"encoding/json"
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

func init() {
	jsonParsers.RegisterParser("indicatorRange", func(arg json.RawMessage) (Condition, error) {
		var params struct {
			Indicator json.RawMessage `json:"indicator"`
			Min       float64         `json:"min"`
			Max       float64         `json:"max"`
		}

		if err := json.Unmarshal(arg, &params); err != nil {
			return nil, fmt.Errorf("failed to parse IndicatorRange parameters: %w", err)
		}

		indicator, err := indicators.FromJSON(params.Indicator)
		if err != nil {
			return nil, fmt.Errorf("failed to parse indicator: %w", err)
		}

		return IndicatorRange(indicator, params.Min, params.Max), nil
	})
}
