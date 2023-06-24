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
	"go.mongodb.org/mongo-driver/mongo"
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

		logger.Debugf("Start query")

		col := database.GetCollection("measures")

		stage1 := bson.D{{"$match", bson.D{{"sensor.sensorId", bson.D{{"$regex", "real"}}}}}}
		stage2 := bson.D{{"$sort", bson.D{{"timestamp", 1}}}}
		stage3 := bson.D{{"$group", bson.D{{"_id", "$sensor.sensorId"}, {"timestamp", bson.D{{"$last", "$timestamp"}}}, {"value", bson.D{{"$last", "$value"}}}}}}
		stage4 := bson.D{{"$sort", bson.D{{"_id", 1}}}}

		cursor, err := col.Aggregate(context.TODO(), mongo.Pipeline{stage1, stage2, stage3, stage4})

		if err != nil {
			logger.Error(err)
			return
		}

		var results []interface{}

		if err = cursor.All(context.TODO(), &results); err != nil {
			logger.Error(err)
		}

		for _, result := range results {
			logger.Infof("Item: %+v\n", result)
		}
	}()
}
