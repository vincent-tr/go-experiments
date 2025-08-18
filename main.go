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
	"time"
)

func gs_main() {
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

func main() {
	dataset, err := backtesting.LoadDataset(
		common.NewMonth(2024, 1),
		common.NewMonth(2024, 12),
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

	fmt.Printf("STRAT JSON: %s\n", modular.ToJSON(builder))

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

	//spew.Dump(metrics)
	printMetricsSummary(metrics)
}
func printMetricsSummary(monthlyMetrics map[common.Month]*backtesting.Metrics) {
	fmt.Printf("\n📊 Trading Summary\n")
	fmt.Printf("==================\n")

	// Aggregate all monthly metrics
	var totalTrades, totalWinningTrades, totalLongTrades, totalShortTrades int
	var totalNetPnL float64
	var totalDuration time.Duration
	var maxDrawdown float64

	// Sort months chronologically
	months := make([]common.Month, 0, len(monthlyMetrics))
	for month := range monthlyMetrics {
		months = append(months, month)
	}

	// Sort months
	for i := 0; i < len(months)-1; i++ {
		for j := i + 1; j < len(months); j++ {
			if months[i].Year() > months[j].Year() ||
				(months[i].Year() == months[j].Year() && months[i].Month() > months[j].Month()) {
				months[i], months[j] = months[j], months[i]
			}
		}
	}

	// Aggregate metrics
	for _, metrics := range monthlyMetrics {
		totalTrades += metrics.TotalTrades
		totalWinningTrades += int(metrics.WinRate * float64(metrics.TotalTrades) / 100)
		totalLongTrades += metrics.LongTrades
		totalShortTrades += metrics.ShortTrades
		totalNetPnL += metrics.NetPnL
		totalDuration += metrics.AvgTradeDuration * time.Duration(metrics.TotalTrades)
		if metrics.MaxDrawdownPct > maxDrawdown {
			maxDrawdown = metrics.MaxDrawdownPct
		}
	}

	// Overall performance
	var profitColor string
	if totalNetPnL > 0 {
		profitColor = "\033[32m" // Green
	} else if totalNetPnL < 0 {
		profitColor = "\033[31m" // Red
	} else {
		profitColor = "\033[37m" // White
	}

	fmt.Printf("💰 Total Profit: %s%.2f\033[0m\n", profitColor, totalNetPnL)
	fmt.Printf("📈 Total Trades: %d\n", totalTrades)
	fmt.Printf("📊 Long Trades: %d\n", totalLongTrades)
	fmt.Printf("📉 Short Trades: %d\n", totalShortTrades)

	overallWinRate := 0.0
	if totalTrades > 0 {
		overallWinRate = float64(totalWinningTrades) / float64(totalTrades) * 100
	}

	var winRateColor string
	if overallWinRate >= 60 {
		winRateColor = "\033[32m" // Green
	} else if overallWinRate >= 40 {
		winRateColor = "\033[33m" // Yellow
	} else {
		winRateColor = "\033[31m" // Red
	}

	fmt.Printf("🎯 Win Rate: %s%.1f%%\033[0m\n", winRateColor, overallWinRate)

	// Calculate profit factor from aggregated data
	if len(monthlyMetrics) > 0 {
		avgProfitFactor := 0.0
		validMonths := 0
		for _, metrics := range monthlyMetrics {
			if metrics.ProfitFactor > 0 {
				avgProfitFactor += metrics.ProfitFactor
				validMonths++
			}
		}
		if validMonths > 0 {
			avgProfitFactor /= float64(validMonths)
			var pfColor string
			if avgProfitFactor >= 1.5 {
				pfColor = "\033[32m" // Green
			} else if avgProfitFactor >= 1.0 {
				pfColor = "\033[33m" // Yellow
			} else {
				pfColor = "\033[31m" // Red
			}
			fmt.Printf("📊 Avg Profit Factor: %s%.2f\033[0m\n", pfColor, avgProfitFactor)
		}
	}

	fmt.Printf("📉 Max Drawdown: \033[31m%.2f%%\033[0m\n", maxDrawdown)

	if totalTrades > 0 {
		avgDuration := totalDuration / time.Duration(totalTrades)
		fmt.Printf("⏱️  Avg Trade Duration: %s\n", avgDuration.String())
	}

	// Monthly breakdown
	if len(monthlyMetrics) > 0 {
		fmt.Printf("\n📅 Monthly Performance\n")
		fmt.Printf("===================\n")

		for _, month := range months {
			metrics := monthlyMetrics[month]
			var monthColor string
			if metrics.NetPnL > 0 {
				monthColor = "\033[32m" // Green
			} else if metrics.NetPnL < 0 {
				monthColor = "\033[31m" // Red
			} else {
				monthColor = "\033[37m" // White
			}

			fmt.Printf("📆 %s: %s%.2f\033[0m (Trades: %d, Win Rate: %.1f%%)\n",
				month.String(), monthColor, metrics.NetPnL, metrics.TotalTrades, metrics.WinRate)
		}
	}

	fmt.Printf("\n")
}
