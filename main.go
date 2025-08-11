package main

import (
	"fmt"
	"go-experiments/brokers/backtesting"
	"go-experiments/common"
	"go-experiments/gridsearch"
	"go-experiments/runner"
	"go-experiments/strategies"
	"go-experiments/traders"
	"go-experiments/traders/modular"
	"go-experiments/traders/modular/indicators"
	"go-experiments/traders/modular/ordercomputer"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	instrument := "EURUSD"

	months := []common.Month{
		common.NewMonth(2023, 1),
		common.NewMonth(2023, 2),
		common.NewMonth(2023, 3),
		common.NewMonth(2023, 4),
		common.NewMonth(2023, 5),
		common.NewMonth(2023, 6),
	}

	runner, err := runner.NewRunner()
	if err != nil {
		panic(err)
	}
	defer runner.Close()

	combos := strategies.BreakoutSpace.GenerateCombinations()

	fmt.Printf("Combined %d strategies\n", len(combos))

	for _, combo := range combos {
		for _, month := range months {
			strategy := buildStrategy(combo)
			if err := runner.SubmitRun(instrument, month, strategy); err != nil {
				panic(err)
			}
		}
	}
}

func buildStrategy(combo gridsearch.Combo) modular.Builder {
	builder := modular.NewBuilder()
	builder.SetHistorySize(250)

	strategies.BreakoutGS(builder.Strategy(), combo)

	builder.RiskManager().SetStopLoss(
		ordercomputer.StopLossATR(indicators.ATR(14), 1.0),
		//ordercomputer.StopLossPipBuffer(3, 15),
	).SetTakeProfit(
		ordercomputer.TakeProfitRatio(2.0),
	)

	builder.CapitalAllocator().SetAllocator(
		ordercomputer.CapitalFixed(10),
	)

	return builder
}

func old_main() {
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

	// ////// TEST
	// raw, err := modular.ToJSON(builder)
	// if err != nil {
	// 	panic(err)
	// }

	// newBuilder, err := modular.FromJSON(raw)
	// if err != nil {
	// 	panic(err)
	// }

	// str1 := modular.Format(builder)
	// str2 := modular.Format(newBuilder)

	// if str1 != str2 {
	// 	fmt.Printf("Original: %s\n", str1)
	// 	fmt.Printf("New: %s\n", str2)
	// 	panic("formatted strings do not match")
	// }

	// ////// TEST

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
