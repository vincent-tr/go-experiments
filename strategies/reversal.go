package strategies

import (
	"go-experiments/common"
	"go-experiments/traders/modular"
	"go-experiments/traders/modular/conditions"
	"go-experiments/traders/modular/indicators"
	"time"
)

func Reversal(strategy modular.StrategyBuilder) {
	strategy.SetFilter(
		conditions.And(
			conditions.HistoryComplete(),
			conditions.NoOpenPositions(),
			conditions.Weekday(time.Tuesday, time.Wednesday, time.Thursday),
			conditions.Session(common.LondonSession),
			conditions.Session(common.NYSession),
			conditions.ExcludeUKHolidays(),
			conditions.ExcludeUSHolidays(),
		),
	)

	strategy.SetLongTrigger(
		conditions.And(
			conditions.PriceThreshold(indicators.EMA(50), conditions.Above),
			conditions.CrossOver(
				indicators.Const(14, 30.0),
				indicators.RSI(14),
				conditions.CrossOverUp,
			),
		),
	)

	strategy.SetShortTrigger(
		conditions.And(
			conditions.PriceThreshold(indicators.EMA(50), conditions.Below),
			conditions.CrossOver(
				indicators.Const(14, 70.0),
				indicators.RSI(14),
				conditions.CrossOverDown,
			),
		),
	)
}
