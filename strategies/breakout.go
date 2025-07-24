package strategies

import (
	"go-experiments/common"
	"go-experiments/traders/modular"
	"go-experiments/traders/modular/conditions"
	"go-experiments/traders/modular/indicators"
	"time"
)

func Breakout(strategy modular.StrategyBuilder) {

	strategy.SetFilter(conditions.And(
		conditions.HistoryComplete(),
		conditions.NoOpenPositions(),

		conditions.Weekday(time.Tuesday, time.Wednesday, time.Thursday),
		conditions.ExcludeUKHolidays(),
		conditions.ExcludeUSHolidays(),
		conditions.Session(common.LondonSession),
		conditions.Session(common.NYSession),

		conditions.IndicatorRange(indicators.RSI(14), 30, 70),
		conditions.Threshold(indicators.ADX(14), 20.0, conditions.Above),
	))

	strategy.SetLongTrigger(
		conditions.And(
			//conditions.PriceThreshold(indicators.EMA(200), conditions.Above),
			conditions.CrossOver(
				indicators.EMA(20),
				indicators.EMA(5),
				conditions.CrossOverUp,
			),
		),
	)

	strategy.SetShortTrigger(
		conditions.And(
			//conditions.PriceThreshold(indicators.EMA(200), conditions.Below),
			conditions.CrossOver(
				indicators.EMA(20),
				indicators.EMA(5),
				conditions.CrossOverDown,
			),
		),
	)
}
