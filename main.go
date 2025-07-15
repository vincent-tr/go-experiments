package main

import (
	"go-experiments/brokers/backtesting"
	"go-experiments/traders"
	"time"
)

func main() {

	brokerConfig := &backtesting.Config{
		BeginDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		Symbol:    "EURUSD",

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

	broker, err := backtesting.NewBroker(brokerConfig)
	if err != nil {
		panic(err)
	}

	traderConfig := &traders.GptConfig{
		EmaFastPeriod:        5,
		EmaSlowPeriod:        20,
		RsiPeriod:            14,
		RsiMin:               30,
		RsiMax:               70,
		StopLossPipBuffer:    3,
		StopLossLookupPeriod: 15,
		TakeProfitRatio:      2.0,
		CapitalRiskPercent:   1.0,
	}

	traders.SetupGptTrader(broker, traderConfig)

	if err := broker.Run(); err != nil {
		panic(err)
	}
}
