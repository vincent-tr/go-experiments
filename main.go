package main

import (
	"go-experiments/brokers/backtesting"
	"go-experiments/common"
	"go-experiments/strategies"
	"go-experiments/traders"
	"go-experiments/traders/modular"
	"go-experiments/traders/modular/indicators"
	"go-experiments/traders/modular/ordercomputer"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	dataset, err := backtesting.LoadDataset(
		common.NewMonth(2023, 4),
		common.NewMonth(2023, 6),
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

		InitialCapital: 100000,
	}

	broker, err := backtesting.NewBroker(brokerConfig, dataset)
	if err != nil {
		panic(err)
	}

	builder := modular.NewBuilder()
	builder.SetHistorySize(250)

	strategies.Breakout(builder.Strategy())

	builder.RiskManager().SetStopLoss(
		ordercomputer.StopLossATR(indicators.ATR(14), 1.0),
		//ordercomputer.StopLossPipBuffer(3, 15),
	).SetTakeProfit(
		ordercomputer.TakeProfitRatio(2.0),
	)

	builder.CapitalAllocator().SetAllocator(
		ordercomputer.CapitalFixed(10),
	)

	if err := traders.SetupModularTrader(broker, builder); err != nil {
		panic(err)
	}
	if err := broker.Run(); err != nil {
		panic(err)
	}

	metrics, err := backtesting.ComputeMetrics(broker)
	if err != nil {
		panic(err)
	}
	spew.Dump(metrics)
}
