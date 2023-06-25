package main

import (
	"context"
	"mylife-tools-server/services/database"
	"mylife-tools-server/services/io/serialization"
	"mylife-tools-server/services/notification"
	"mylife-tools-server/services/sessions"
	"mylife-tools-server/services/store"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

var measures *store.Container[*Measure]

type Measure struct {
	id        string
	timestamp time.Time
	value     float64
	sensor    string // or *Sensor ?
}

func (measure *Measure) Id() string {
	return measure.id
}

func (measure *Measure) Marshal() (interface{}, error) {
	helper := serialization.NewStructMarshallerHelper()

	helper.Add("_id", measure.id)
	helper.Add("timestamp", measure.timestamp)
	helper.Add("value", measure.value)
	helper.Add("sensor", measure.sensor)

	return helper.Build()
}

func notifyMeasures(session *sessions.Session, arg struct{}) (uint64, error) {
	viewId := notification.NotifyView[*Measure](session, measures)
	return viewId, nil
}

func initMeasures() {

	measures = store.NewContainer[*Measure]("measures")

	measures.Set(&Measure{
		id:        "id1",
		timestamp: time.Now(),
		value:     42,
		sensor:    "the sensor",
	})

	go func() {
		time.Sleep(time.Second)

		for {
			begin := time.Now()
			fetchResults()
			elapsed := time.Since(begin).Seconds() * 1000
			logger.WithField("elapsedMs", elapsed).Debug("Fetch results")

			time.Sleep(time.Second * 5)
		}
	}()

}

func fetchResults() {
	col := database.GetCollection("measures")

	cursor, err := col.Aggregate(context.TODO(), []bson.M{
		{"$sort": bson.M{"sensor.sensorId": 1, "timestamp": -1}},
		{"$group": bson.M{"_id": "$sensor.sensorId", "timestamp": bson.M{"$first": "$timestamp"}, "value": bson.M{"$first": "$value"}}},
	})

	if err != nil {
		logger.Error(err)
		return
	}

	defer cursor.Close(context.TODO())

	var results []interface{}

	if err = cursor.All(context.TODO(), &results); err != nil {
		logger.Error(err)
	}

	for _, result := range results {
		logger.Infof("Item: %+v\n", result)
	}
}
