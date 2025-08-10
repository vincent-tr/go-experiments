package runner

import (
	"fmt"
	"go-experiments/brokers/backtesting"
	"go-experiments/common"
	"go-experiments/traders"
	"go-experiments/traders/modular"
)

var log = common.NewLogger("runner")

type Runner struct {
	db       *Database
	datasets *datasets
	pool     *TaskPool
}

func NewRunner() (*Runner, error) {
	db, err := OpenDatabase()
	if err != nil {
		return nil, err
	}

	return &Runner{
		db:       db,
		datasets: newDatasets(),
		pool:     NewTaskPool(),
	}, nil
}

func (r *Runner) Close() {
	r.db.Close()
	r.pool.Close()
}

func (r *Runner) SubmitRun(instrument string, month common.Month, strategy modular.Builder) error {
	// Try to see if output is already cached
	strategyStr := modular.ToJSON(strategy)

	run, err := r.db.FindRun(instrument, month.String(), strategyStr)
	if err != nil {
		return err
	}

	if run != nil {
		log.Info("Run already exists for %s %s: %s", instrument, month.String(), strategy.Format().Compact())
		return nil
	}

	// enqueue run
	r.pool.Submit(func() {
		if err := r.run(instrument, month, strategy); err != nil {
			log.Error("Failed to run strategy for %s %s: %v", instrument, month.String(), err)
		}
	})

	return nil
}

func (r *Runner) run(instrument string, month common.Month, strategy modular.Builder) error {
	log.Info("Running strategy for %s %s: %s", instrument, month.String(), strategy.Format().Compact())

	dataset, err := r.datasets.Get(instrument, month)
	if err != nil {
		return fmt.Errorf("failed to get dataset for %s %s: %w", instrument, month.String(), err)
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
		return fmt.Errorf("failed to create broker: %w", err)
	}

	if err := traders.SetupModularTrader(broker, strategy); err != nil {
		return fmt.Errorf("failed to setup trader: %w", err)
	}
	if err := broker.Run(); err != nil {
		return fmt.Errorf("failed to run broker: %w", err)
	}

	metrics, err := backtesting.ComputeMetrics(broker)
	if err != nil {
		return fmt.Errorf("failed to compute metrics: %w", err)
	}

	var metrics0 *backtesting.Metrics

	switch len(metrics) {
	case 0:
		// No metrics means no position has been taken
		metrics0 = &backtesting.Metrics{}

	case 1:
		var ok bool
		metrics0, ok = metrics[month]
		if !ok {
			return fmt.Errorf("no metrics found for month %s", month.String())
		}

	default:
		return fmt.Errorf("expected exactly one metric, got %d", len(metrics))
	}

	if err := r.saveResult(instrument, month, strategy, metrics0); err != nil {
		return fmt.Errorf("failed to save result: %w", err)
	}

	log.Info("Run completed for %s %s: %s", instrument, month.String(), strategy.Format().Compact())
	return nil
}

func (r *Runner) saveResult(instrument string, month common.Month, strategy modular.Builder, metrics *backtesting.Metrics) error {
	strategyStr := modular.ToJSON(strategy)

	run := &run{
		Key:        r.db.ComputeKey(instrument, month.String(), strategyStr),
		Instrument: instrument,
		TimeRange:  month.String(),
		Strategy:   strategyStr,
		Metrics:    *metrics,
	}

	if err := r.db.SaveRun(run); err != nil {
		return fmt.Errorf("failed to save run: %w", err)
	}

	return nil
}
