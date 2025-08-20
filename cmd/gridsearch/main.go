package main

import (
	"fmt"
	"go-experiments/common"
	"go-experiments/gridsearch"
	"go-experiments/runner"
	"go-experiments/strategies"
	"go-experiments/traders/modular"
	"go-experiments/traders/modular/indicators"
	"go-experiments/traders/modular/ordercomputer"
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
