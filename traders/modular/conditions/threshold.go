package conditions

import (
	"encoding/json"
	"fmt"
	"go-experiments/traders/modular/context"
	"go-experiments/traders/modular/formatter"
	"go-experiments/traders/modular/indicators"
)

type Direction int

const (
	Above Direction = iota
	Below
)

func (d Direction) String() string {
	switch d {
	case Above:
		return "Above"
	case Below:
		return "Below"
	default:
		return "Unknown"
	}
}

// Threshold checks if the value of an indicator is above or below a specified threshold.
func Threshold(indicator indicators.Indicator, threshold float64, direction Direction) Condition {
	return newCondition(
		func(ctx context.TraderContext) bool {
			values := indicator.Values(ctx)
			if len(values) == 0 {
				return false
			}
			value := values[len(values)-1]

			switch direction {
			case Above:
				return value >= threshold
			case Below:
				return value <= threshold
			default:
				panic(fmt.Sprintf("unknown threshold direction: %d", direction))
			}
		},
		func() *formatter.FormatterNode {
			return formatter.Format("Threshold",
				indicator.Format(),
				formatter.Format(fmt.Sprintf("Value: %.2f", threshold)),
				formatter.Format(fmt.Sprintf("Direction: %s", direction.String())),
			)
		},
	)
}

func init() {
	jsonParsers.RegisterParser("threshold", func(arg json.RawMessage) (Condition, error) {
		var params struct {
			Indicator json.RawMessage `json:"indicator"`
			Threshold float64         `json:"threshold"`
			Direction string          `json:"direction"`
		}

		if err := json.Unmarshal(arg, &params); err != nil {
			return nil, fmt.Errorf("failed to parse Threshold parameters: %w", err)
		}

		indicator, err := indicators.FromJSON(params.Indicator)
		if err != nil {
			return nil, fmt.Errorf("failed to parse indicator: %w", err)
		}

		var direction Direction
		switch params.Direction {
		case "above":
			direction = Above
		case "below":
			direction = Below
		default:
			return nil, fmt.Errorf("unknown direction: %s", params.Direction)
		}

		return Threshold(indicator, params.Threshold, direction), nil
	})
}

// PriceThreshold checks if the current price is above or below the value of an indicator.
func PriceThreshold(indicator indicators.Indicator, direction Direction) Condition {
	return newCondition(
		func(ctx context.TraderContext) bool {
			values := indicator.Values(ctx)
			if len(values) == 0 {
				return false
			}
			value := values[len(values)-1]
			entryPrice := ctx.EntryPrice()

			switch direction {
			case Above:
				return entryPrice >= value
			case Below:
				return entryPrice <= value
			default:
				panic(fmt.Sprintf("unknown threshold direction: %d", direction))
			}
		},
		func() *formatter.FormatterNode {
			return formatter.Format("PriceThreshold",
				indicator.Format(),
				formatter.Format(fmt.Sprintf("Direction: %s", direction.String())),
			)
		},
	)
}

func init() {
	jsonParsers.RegisterParser("priceThreshold", func(arg json.RawMessage) (Condition, error) {
		var params struct {
			Indicator json.RawMessage `json:"indicator"`
			Direction string          `json:"direction"`
		}

		if err := json.Unmarshal(arg, &params); err != nil {
			return nil, fmt.Errorf("failed to parse PriceThreshold parameters: %w", err)
		}

		indicator, err := indicators.FromJSON(params.Indicator)
		if err != nil {
			return nil, fmt.Errorf("failed to parse indicator: %w", err)
		}

		var direction Direction
		switch params.Direction {
		case "above":
			direction = Above
		case "below":
			direction = Below
		default:
			return nil, fmt.Errorf("unknown direction: %s", params.Direction)
		}

		return PriceThreshold(indicator, direction), nil
	})
}
