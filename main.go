package main

import (
	"go-experiments/brokers/backtesting"
	"time"
)

func main() {
	symbol := "EURUSD"
	beginDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	backtesting.NewBroker(beginDate, endDate, symbol)

}
