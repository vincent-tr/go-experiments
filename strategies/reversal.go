package strategies

import (
	"go-experiments/common"
	"go-experiments/traders/modular"
	"go-experiments/traders/modular/condition"
	"go-experiments/traders/modular/indicators"
	"time"
)

func Reversal(strategy modular.StrategyBuilder) {
	strategy.SetFilter(
		condition.And(
			condition.HistoryComplete(),
			condition.NoOpenPositions(),
			condition.Weekday(time.Tuesday, time.Wednesday, time.Thursday),
			condition.Session(common.LondonSession),
			condition.Session(common.NYSession),
			condition.ExcludeUKHolidays(),
			condition.ExcludeUSHolidays(),
		),
	)

	strategy.SetLongTrigger(
		condition.CrossOver(
			indicators.Const(14, 30.0),
			indicators.RSI(14),
			condition.CrossOverUp,
		),
	)

	strategy.SetShortTrigger(
		condition.CrossOver(
			indicators.Const(14, 70.0),
			indicators.RSI(14),
			condition.CrossOverDown,
		),
	)
}
