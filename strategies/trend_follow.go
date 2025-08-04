package strategies

import (
	"go-experiments/common"
	"go-experiments/traders/modular"
	"go-experiments/traders/modular/conditions"
	"go-experiments/traders/modular/indicators"
	"time"
)

func TrendFollow(strategy modular.StrategyBuilder) {
	strategy.SetFilter(
		conditions.And(
			conditions.HistoryUsable(),
			conditions.NoOpenPositions(),

			conditions.Weekday(time.Tuesday, time.Wednesday, time.Thursday),
			conditions.ExcludeUKHolidays(),
			conditions.ExcludeUSHolidays(),
			conditions.Session(common.LondonSession),
			conditions.Session(common.NYSession),

			conditions.Threshold(indicators.ADX(14), 20.0, conditions.Above),
		),
	)

	strategy.SetLongTrigger(
		conditions.And(
			conditions.IndicatorRange(indicators.RSI(14), 40, 70),
			conditions.PriceThreshold(indicators.EMA(50), conditions.Above),
			conditions.CrossOver(
				indicators.EMA(50),
				indicators.EMA(20),
				conditions.CrossOverUp,
			),
		),
	)

	strategy.SetShortTrigger(
		conditions.And(
			conditions.IndicatorRange(indicators.RSI(14), 30, 60),
			conditions.PriceThreshold(indicators.EMA(50), conditions.Below),
			conditions.CrossOver(
				indicators.EMA(50),
				indicators.EMA(20),
				conditions.CrossOverDown,
			),
		),
	)
}
