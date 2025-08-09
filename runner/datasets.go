package runner

import (
	"fmt"
	"go-experiments/brokers/backtesting"
	"go-experiments/common"
	"sync"
)

type datasets struct {
	datasets map[string]*backtesting.Dataset
	lock     sync.Mutex
}

func newDatasets() *datasets {
	return &datasets{
		datasets: make(map[string]*backtesting.Dataset),
		lock:     sync.Mutex{},
	}
}

func (d *datasets) Get(instrument string, month common.Month) (*backtesting.Dataset, error) {
	d.lock.Lock()
	defer d.lock.Unlock()

	key := fmt.Sprintf("%s-%s", instrument, month.String())

	if dataset, exists := d.datasets[key]; exists {
		return dataset, nil
	}

	dataset, err := backtesting.LoadDataset(month, month, instrument)
	if err != nil {
		return nil, err
	}

	d.datasets[key] = dataset
	return dataset, nil
}
