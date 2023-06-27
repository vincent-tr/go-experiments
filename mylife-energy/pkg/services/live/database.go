package live

import (
	"context"
	"mylife-energy/pkg/entities"
	"mylife-tools-server/services/database"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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

func fetchResults() ([]mongoMeasure, error) {
	col := database.GetCollection("measures")

	cursor, err := col.Aggregate(context.TODO(), []bson.M{
		{"$sort": bson.M{"sensor.sensorId": 1, "timestamp": -1}},
		{"$group": bson.M{"_id": "$sensor.sensorId", "timestamp": bson.M{"$first": "$timestamp"}, "value": bson.M{"$first": "$value"}, "sensor": bson.M{"$first": "$sensor"}}},
	})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.TODO())

	var results []mongoMeasure

	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

func makeMeasureFromData(data *mongoMeasure) *entities.Measure {
	return entities.NewMeasure(entities.MeasureData{
		Id:        data.Id,
		Sensor:    data.Sensor.SensorId,
		Timestamp: data.Timestamp,
		Value:     data.Value,
	})
}

func makeSensorFromData(data *mongoSensor) *entities.Sensor {
	return entities.NewSensor(entities.SensorData{
		Id:                data.SensorId,
		DeviceClass:       data.DeviceClass,
		StateClass:        data.StateClass,
		UnitOfMeasurement: data.UnitOfMeasurement,
		AccuracyDecimals:  data.AccuracyDecimals,
	})
}
