package main

import (
	"go-experiments/brokers/backtesting"
	"go-experiments/common"
	"go-experiments/traders"
	"go-experiments/traders/modular"
	"go-experiments/traders/modular/condition"
	"go-experiments/traders/modular/indicators"
	"go-experiments/traders/modular/ordercomputer"
	"time"
)

func main() {
	dataset, err := backtesting.LoadDataset(
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		"EURUSD",
	)

	if err != nil {
		panic(err)
	}

	brokerConfig := &backtesting.Config{
		// For backtesting, we assume a lot size of 1 for simplicity.
		// In a real broker, this would be the number of units per lot.
		// Not that using IG broker, EUR/USD Mini has also a size of 1.
		LotSize: 1,

		// Leverage is the ratio of the amount of capital that a trader must put up to open a position.
		// For example, if the leverage is 30, it means that for every 1 unit of capital,
		// the trader can control 30 units of the asset.
		// This is a common leverage ratio in forex trading.
		Leverage: 30.0,

		InitialCapital: 1000,
	}

	broker, err := backtesting.NewBroker(brokerConfig, dataset)
	if err != nil {
		panic(err)
	}
	/*
		traderConfig := &traders.GptConfig{
			HistorySize:           100,
			EmaFastPeriod:         5,
			EmaSlowPeriod:         20,
			RsiPeriod:             14,
			RsiMin:                30,
			RsiMax:                70,
			StopLossAtrEnabled:    true,
			StopLossAtrPeriod:     14,
			StopLossAtrMultiplier: 1,
			//StopLossPipBuffer:     3,
			//StopLossLookupPeriod:  15,
			TakeProfitRatio:    2.0,
			CapitalRiskPercent: 1.0,
			AdxEnabled:         true,
			AdxPeriod:          14,
			AdxThreshold:       20.0,
		}

		traders.SetupGptTrader(broker, traderConfig)
	*/
	builder := modular.NewBuilder()
	builder.SetHistorySize(100)

	builder.Strategy().SetFilter(condition.And(
		condition.HistoryComplete(),
		condition.NoOpenPositions(),

		condition.Weekday(time.Tuesday, time.Wednesday, time.Thursday),
		condition.ExcludeUKHolidays(),
		condition.ExcludeUSHolidays(),
		condition.Session(common.LondonSession),
		condition.Session(common.NYSession),

		condition.IndicatorRange(indicators.RSI(14), 30, 70),
		condition.Threshold(indicators.ADX(14), 20.0),
	))

	builder.Strategy().SetLongTrigger(
		condition.CrossOver(
			indicators.EMA(20),
			indicators.EMA(5),
			condition.CrossOverUp,
		),
	)

	builder.Strategy().SetShortTrigger(
		condition.CrossOver(
			indicators.EMA(20),
			indicators.EMA(5),
			condition.CrossOverDown,
		),
	)

	builder.RiskManager().SetStopLoss(
		ordercomputer.StopLossATR(indicators.ATR(14), 1.0),
		//ordercomputer.StopLossPipBuffer(3, 15),
	).SetTakeProfit(
		ordercomputer.TakeProfitRatio(2.0),
	)

	builder.CapitalAllocator().SetAllocator(
		ordercomputer.CapitalRiskPercent(1.0),
	)

	if err := traders.SetupModularTrader(broker, builder); err != nil {
		panic(err)
	}
	if err := broker.Run(); err != nil {
		panic(err)
	}
}
