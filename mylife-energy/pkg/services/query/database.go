package query

import (
	"context"
	"mylife-energy/pkg/entities"
	"mylife-tools-server/services/database"
	"mylife-tools-server/utils"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type mongoMeasure struct {
	Id        string      `bson:"_id"`
	Timestamp time.Time   `bson:"timestamp"`
	Value     float64     `bson:"value"`
	Sensor    mongoSensor `bson:"sensor"`
}

type mongoSensor struct {
	SensorId          string `bson:"sensorId"`
	DeviceClass       string `bson:"deviceClass"`
	StateClass        string `bson:"stateClass"`
	UnitOfMeasurement string `bson:"unitOfMeasurement"`
	AccuracyDecimals  uint   `bson:"accuracyDecimals"`
}

func dbFetch(ctx context.Context, cursorBuilder func(ctx context.Context, col *mongo.Collection) (*mongo.Cursor, error)) ([]mongoMeasure, error) {
	col := database.GetCollection("measures")

	logger.Debug("Query begin")
	tmr := utils.NewTimer()
	defer func() {
		logger.WithField("elapsedMs", tmr.ElapsedMs()).Debug("Query end")
	}()

	cursor, err := cursorBuilder(ctx, col)

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var results []mongoMeasure

	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func makeMeasureFromData(data *mongoMeasure) *entities.Measure {
	return entities.NewMeasure(&entities.MeasureData{
		Id:        data.Id,
		Sensor:    data.Sensor.SensorId,
		Timestamp: data.Timestamp,
		Value:     data.Value,
	})
}

func makeSensorFromData(data *mongoSensor) *entities.Sensor {
	return entities.NewSensor(&entities.SensorData{
		Id:                data.SensorId,
		DeviceClass:       data.DeviceClass,
		StateClass:        data.StateClass,
		UnitOfMeasurement: data.UnitOfMeasurement,
		AccuracyDecimals:  data.AccuracyDecimals,
	})
}
