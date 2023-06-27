package query

import (
	"context"
	"mylife-energy/pkg/entities"
	"mylife-tools-server/log"
	"mylife-tools-server/services"

	"go.mongodb.org/mongo-driver/mongo"
)

var logger = log.CreateLogger("mylife:energy:query")

type Result struct {
	Measure *entities.Measure
	Sensor  *entities.Sensor
}

type queryService struct {
}

func (service *queryService) Init(arg interface{}) error {
	return nil
}

func (service *queryService) Terminate() error {
	return nil
}

func (service *queryService) ServiceName() string {
	return "query"
}

func (service *queryService) Dependencies() []string {
	// io because we use io queue
	return []string{"database"}
}

func init() {
	services.Register(&queryService{})
}

func (service *queryService) exec(ctx context.Context, cursorBuilder func(ctx context.Context, col *mongo.Collection) (*mongo.Cursor, error)) ([]Result, error) {
	mongoResults, err := dbFetch(ctx, cursorBuilder)
	if err != nil {
		return nil, err
	}

	var results = make([]Result, 0, len(mongoResults))

	for _, mongoResult := range mongoResults {
		results = append(results, Result{
			Measure: makeMeasureFromData(&mongoResult),
			Sensor:  makeSensorFromData(&mongoResult.Sensor),
		})
	}

	return results, nil
}

func getService() *queryService {
	return services.GetService[*queryService]("query")
}

// Public access

func Aggregate(ctx context.Context, pipeline interface{}) ([]Result, error) {
	return getService().exec(ctx, func(ctx context.Context, col *mongo.Collection) (*mongo.Cursor, error) {
		return col.Aggregate(ctx, pipeline)
	})
}

func Find(ctx context.Context, filter interface{}) ([]Result, error) {
	return getService().exec(ctx, func(ctx context.Context, col *mongo.Collection) (*mongo.Cursor, error) {
		return col.Find(ctx, filter)
	})
}
