package main

import (
	"go-experiments/brokers/backtesting"
	"go-experiments/traders"
	"time"
)

func main() {
	symbol := "EURUSD"
	//beginDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	//endDate := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

	// For backtesting, we assume a lot size of 1 for simplicity.
	// In a real broker, this would be the number of units per lot.
	// Not that using IG broker, EUR/USD Mini has also a size of 1.
	lotSize := 1

	// Leverage is the ratio of the amount of capital that a trader must put up to open a position.
	// For example, if the leverage is 30, it means that for every 1 unit of capital,
	// the trader can control 30 units of the asset.
	// This is a common leverage ratio in forex trading.
	leverage := 30.0

	beginDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	broker, err := backtesting.NewBroker(beginDate, endDate, symbol, lotSize, leverage, 1000)
	if err != nil {
		panic(err)
	}

	traders.SetupGptTrader(broker)

	if err := broker.Run(); err != nil {
		panic(err)
	}
}
