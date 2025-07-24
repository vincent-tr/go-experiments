package condition

import (
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
