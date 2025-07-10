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

	beginDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC)
	broker, err := backtesting.NewBroker(beginDate, endDate, symbol, 10000)
	if err != nil {
		panic(err)
	}

	traders.SetupBasicTrader(broker)

	if err := broker.Run(); err != nil {
		panic(err)
	}
}
