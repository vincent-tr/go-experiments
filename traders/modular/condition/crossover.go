package condition

import (
	"go-experiments/traders/modular/context"
	"go-experiments/traders/modular/formatter"
	"go-experiments/traders/modular/indicators"
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
	)
}
