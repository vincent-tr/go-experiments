package runner

import (
	"go-experiments/brokers/backtesting"
	"go-experiments/traders/modular"
)

type Runner struct {
	db       *Database
	datasets map[string]*backtesting.Dataset
	pool     *TaskPool
}

func NewRunner() (*Runner, error) {
	db, err := OpenDatabase()
	if err != nil {
		return nil, err
	}

	return &Runner{
		db:       db,
		datasets: make(map[string]*backtesting.Dataset),
		pool:     NewTaskPool(),
	}, nil
}

func (r *Runner) Close() {
	r.db.Close()
	r.pool.Close()
}

func (r *Runner) SubmitRun(instrument string, year int, month int, strategy modular.Builder) error {
	// TODO
	panic("not implemented")
}
