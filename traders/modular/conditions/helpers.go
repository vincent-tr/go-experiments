package conditions

import (
	"encoding/json"
	"go-experiments/traders/modular/context"
	"go-experiments/traders/modular/formatter"
)

func And(conditions ...Condition) Condition {
	return newCondition(
		func(ctx context.TraderContext) bool {
			for _, condition := range conditions {
				if !condition.Execute(ctx) {
					return false
				}
			}
			return true
		},
		func() *formatter.FormatterNode {
			return formatter.FormatWithChildren("And", conditions...)
		},
	)
}

func init() {
	jsonParsers.RegisterParser("and", func(arg json.RawMessage) (Condition, error) {
		var conditions []json.RawMessage
		if err := json.Unmarshal(arg, &conditions); err != nil {
			return nil, err
		}

		var parsedConditions []Condition
		for _, cond := range conditions {
			condition, err := FromJSON(cond)
			if err != nil {
				return nil, err
			}
			parsedConditions = append(parsedConditions, condition)
		}

		return And(parsedConditions...), nil
	})
}

func Or(conditions ...Condition) Condition {
	return newCondition(
		func(ctx context.TraderContext) bool {
			for _, condition := range conditions {
				if condition.Execute(ctx) {
					return true
				}
			}
			return false
		},
		func() *formatter.FormatterNode {
			return formatter.FormatWithChildren("Or", conditions...)
		},
	)
}

func init() {
	jsonParsers.RegisterParser("or", func(arg json.RawMessage) (Condition, error) {
		var conditions []json.RawMessage
		if err := json.Unmarshal(arg, &conditions); err != nil {
			return nil, err
		}

		var parsedConditions []Condition
		for _, cond := range conditions {
			condition, err := FromJSON(cond)
			if err != nil {
				return nil, err
			}
			parsedConditions = append(parsedConditions, condition)
		}

		return Or(parsedConditions...), nil
	})
}

func HistoryUsable() Condition {
	return newCondition(
		func(ctx context.TraderContext) bool {
			return ctx.HistoricalData().IsUsable()
		},
		func() *formatter.FormatterNode {
			return formatter.Format("HistoryUsable")
		},
	)
}

func init() {
	jsonParsers.RegisterParser("historyUsable", func(arg json.RawMessage) (Condition, error) {
		return HistoryUsable(), nil
	})
}

func NoOpenPositions() Condition {
	return newCondition(
		func(ctx context.TraderContext) bool {
			return len(ctx.OpenPositions()) == 0
		},
		func() *formatter.FormatterNode {
			return formatter.Format("NoOpenPositions")
		},
	)
}

func init() {
	jsonParsers.RegisterParser("noOpenPositions", func(arg json.RawMessage) (Condition, error) {
		return NoOpenPositions(), nil
	})
}
