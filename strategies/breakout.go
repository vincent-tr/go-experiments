package strategies

import (
	"go-experiments/common"
	"go-experiments/gridsearch"
	"go-experiments/traders/modular"
	"go-experiments/traders/modular/conditions"
	"go-experiments/traders/modular/indicators"
	"time"
)

func Breakout(strategy modular.StrategyBuilder) {

	strategy.SetFilter(conditions.And(
		conditions.HistoryUsable(),
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

// Breakout strategy parameters to tune via grid search:
//
// 1. Indicator parameters
// - RSI period (e.g., 7, 14, 21)
// - RSI lower bound (e.g., 25, 30, 35)
// - RSI upper bound (e.g., 65, 70, 75)
// - ADX period (e.g., 7, 14, 21)
// - ADX threshold (e.g., 15.0, 20.0, 25.0)
// - EMA short period (e.g., 5, 8, 10)
// - EMA long period (e.g., 20, 30, 50)
//
// 2. Session / day filters
// - Days of week to trade (e.g., Tue-Thu vs Mon-Fri)
// - Session selection (London only, NY only, both)
//
// 3. Risk management parameters (for later)
// - StopLoss ATR period
// - StopLoss ATR multiplier
// - TakeProfit ratio
//
// 4. Optional trend filter
// - EMA200 trend filter: enabled or disabled

var BreakoutSpace = gridsearch.ParameterSpace{
	"RSIPeriod":      {7, 14, 21},
	"RSILower":       {25.0, 30.0, 35.0},
	"RSIUpper":       {65.0, 70.0, 75.0},
	"ADXPeriod":      {7, 14, 21},
	"ADXThreshold":   {15.0, 20.0, 25.0},
	"ShortEMAPeriod": {5, 8, 10},
	"LongEMAPeriod":  {20, 30, 50},
	"TradeDays":      {"TueThu", "MonFri"},
	// "Session":        {"London", "NewYork", "Both"},
	"TrendFilter": {true, false},
}

func BreakoutGS(strategy modular.StrategyBuilder, c gridsearch.Combo) {

	var tradeDays []time.Weekday
	switch c.String("TradeDays") {
	case "TueThu":
		tradeDays = []time.Weekday{time.Tuesday, time.Wednesday, time.Thursday}
	case "MonFri":
		tradeDays = []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday}
	}

	strategy.SetFilter(conditions.And(
		conditions.HistoryUsable(),
		conditions.NoOpenPositions(),

		conditions.Weekday(tradeDays...),
		conditions.ExcludeUKHolidays(),
		conditions.ExcludeUSHolidays(),
		conditions.Session(common.LondonSession),
		conditions.Session(common.NYSession),

		conditions.IndicatorRange(indicators.RSI(c.Int("RSIPeriod")), c.Float("RSILower"), c.Float("RSIUpper")),
		conditions.Threshold(indicators.ADX(c.Int("ADXPeriod")), c.Float("ADXThreshold"), conditions.Above),
	))

	var trendLong conditions.Condition
	var trendShort conditions.Condition

	if c.Bool("TrendFilter") {
		trendLong = conditions.PriceThreshold(indicators.EMA(200), conditions.Above)
		trendShort = conditions.PriceThreshold(indicators.EMA(200), conditions.Below)
	} else {
		trendLong = conditions.True()
		trendShort = conditions.True()
	}

	strategy.SetLongTrigger(
		conditions.And(
			trendLong,
			conditions.CrossOver(
				indicators.EMA(c.Int("LongEMAPeriod")),
				indicators.EMA(c.Int("ShortEMAPeriod")),
				conditions.CrossOverUp,
			),
		),
	)

	strategy.SetShortTrigger(
		conditions.And(
			trendShort,
			conditions.CrossOver(
				indicators.EMA(c.Int("LongEMAPeriod")),
				indicators.EMA(c.Int("ShortEMAPeriod")),
				conditions.CrossOverDown,
			),
		),
	)
}
