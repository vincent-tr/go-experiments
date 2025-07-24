package strategies

import (
	"go-experiments/common"
	"go-experiments/traders/modular"
	"go-experiments/traders/modular/condition"
	"go-experiments/traders/modular/indicators"
	"time"
)

func Breakout(strategy modular.StrategyBuilder) {

	strategy.SetFilter(condition.And(
		condition.HistoryComplete(),
		condition.NoOpenPositions(),

		condition.Weekday(time.Tuesday, time.Wednesday, time.Thursday),
		condition.ExcludeUKHolidays(),
		condition.ExcludeUSHolidays(),
		condition.Session(common.LondonSession),
		condition.Session(common.NYSession),

		condition.IndicatorRange(indicators.RSI(14), 30, 70),
		condition.Threshold(indicators.ADX(14), 20.0, condition.Above),
	))

	strategy.SetLongTrigger(
		condition.And(
			//condition.PriceThreshold(indicators.EMA(200), condition.Above),
			condition.CrossOver(
				indicators.EMA(20),
				indicators.EMA(5),
				condition.CrossOverUp,
			),
		),
	)

	strategy.SetShortTrigger(
		condition.And(
			//condition.PriceThreshold(indicators.EMA(200), condition.Below),
			condition.CrossOver(
				indicators.EMA(20),
				indicators.EMA(5),
				condition.CrossOverDown,
			),
		),
	)
}
