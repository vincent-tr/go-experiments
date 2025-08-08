package conditions

import (
	"encoding/json"
	"fmt"
	"go-experiments/traders/modular/context"
	"go-experiments/traders/modular/formatter"
	"go-experiments/traders/modular/indicators"
	"go-experiments/traders/modular/marshal"
)

type CrossOverDirection int

const (
	CrossOverUp CrossOverDirection = iota
	CrossOverDown
)

func CrossOver(reference, test indicators.Indicator, direction CrossOverDirection) Condition {
	return newCondition(
		func(ctx context.TraderContext) bool {
			refs := reference.Values(ctx)
			tests := test.Values(ctx)
			if len(refs) < 2 || len(tests) < 2 {
				return false
			}

			currRef := refs[len(refs)-1]
			currTest := tests[len(tests)-1]
			prevRef := refs[len(refs)-2]
			prevTest := tests[len(tests)-2]

			switch direction {
			case CrossOverUp:
				return prevTest < prevRef && currTest > currRef
			case CrossOverDown:
				return prevTest > prevRef && currTest < currRef
			default:
				panic("unknown crossover direction")
			}
		},
		func() *formatter.FormatterNode {
			return formatter.Format("CrossOver",
				formatter.FormatWithChildren("Reference", reference),
				formatter.FormatWithChildren("Test", test),
			)
		},
		func() (string, any) {
			var directionStr string
			switch direction {
			case CrossOverUp:
				directionStr = "up"
			case CrossOverDown:
				directionStr = "down"
			default:
				panic("unknown crossover direction")
			}

			return "crossover", map[string]any{
				"reference": marshal.ToJSON(reference),
				"test":      marshal.ToJSON(test),
				"direction": directionStr,
			}
		},
	)
}

func init() {
	jsonParsers.RegisterParser("crossover", func(arg json.RawMessage) (Condition, error) {
		var params struct {
			Reference json.RawMessage `json:"reference"`
			Test      json.RawMessage `json:"test"`
			Direction string          `json:"direction"`
		}

		if err := json.Unmarshal(arg, &params); err != nil {
			return nil, err
		}

		reference, err := indicators.FromJSON(params.Reference)
		if err != nil {
			return nil, fmt.Errorf("failed to parse reference indicator: %w", err)
		}

		test, err := indicators.FromJSON(params.Test)
		if err != nil {
			return nil, fmt.Errorf("failed to parse test indicator: %w", err)
		}

		var direction CrossOverDirection
		switch params.Direction {
		case "up":
			direction = CrossOverUp
		case "down":
			direction = CrossOverDown
		default:
			return nil, fmt.Errorf("invalid crossover direction: %s", params.Direction)
		}

		return CrossOver(reference, test, direction), nil
	})
}
