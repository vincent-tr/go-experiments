package live

import (
	"context"
	"mylife-energy/pkg/entities"
	"mylife-energy/pkg/services/query"
	"mylife-tools-server/services/io"
	"mylife-tools-server/services/store"
	"mylife-tools-server/utils"
	"sync"
	"time"
)

type fetcher struct {
	worker      *utils.Worker
	dbContext   context.Context
	dbTerminate context.CancelFunc
	measures    *store.Container[*entities.Measure]
	sensors     *store.Container[*entities.Sensor]
	pendingSync *sync.WaitGroup
}

func makeFetcher() *fetcher {
	f := &fetcher{
		measures:    store.NewContainer[*entities.Measure]("measures"),
		sensors:     store.NewContainer[*entities.Sensor]("sensors"),
		pendingSync: &sync.WaitGroup{},
	}

	f.dbContext, f.dbTerminate = context.WithCancel(context.Background())

	f.worker = utils.NewWorker(f.workerEntry)

	return f
}

func (f *fetcher) terminate() {
	f.dbTerminate()
	f.worker.Terminate()
	f.pendingSync.Wait()

	f.measures = nil
	f.sensors = nil
}

func (f *fetcher) workerEntry(exit chan struct{}) {

	for {
		select {
		case <-exit:
			return
		case <-time.After(10 * time.Second):
			f.sync()
		}
	}
}

func (f *fetcher) sync() {
	builder := query.NewPipelineBuilder().
		Sort(
			query.SortField{Name: "sensor.sensorId", Order: query.Asc},
			query.SortField{Name: "timestamp", Order: query.Desc},
		).
		Group(
			"$sensor.sensorId",
			query.GroupField{Name: "timestamp", Accumulator: "$first", Expression: "$timestamp"},
			query.GroupField{Name: "value", Accumulator: "$first", Expression: "$value"},
			query.GroupField{Name: "sensor", Accumulator: "$first", Expression: "$sensor"},
		)

	results, err := query.Aggregate(f.dbContext, builder.Build())

	if err != nil {
		logger.WithError(err).Error("Error fetching results")
		return
	}

	f.pendingSync.Add(1)

	err = io.SubmitIoTask("live/fetch", func() {
		newMeasures := make([]*entities.Measure, 0, len(results))
		newSensors := make([]*entities.Sensor, 0, len(results))

		for _, result := range results {
			newMeasures = append(newMeasures, result.Measure)
			newSensors = append(newSensors, result.Sensor)
		}

		f.measures.ReplaceAll(newMeasures, entities.MeasuresEqual)
		f.sensors.ReplaceAll(newSensors, entities.SensorsEqual)

		f.pendingSync.Done()
	})

	if err != nil {
		f.pendingSync.Done()

		logger.WithError(err).Error("Error submitting io task")
		return
	}
}
