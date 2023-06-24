package main

import (
	"mylife-tools-server/services/io/serialization"
	"mylife-tools-server/services/notification"
	"mylife-tools-server/services/sessions"
	"mylife-tools-server/services/store"
)

var sensors *store.Container[*Sensor]

type Sensor struct {
	id                string
	sensorId          string
	deviceClass       string
	stateClass        string
	unitOfMeasurement string
	accuracyDecimals  uint
}

func (sensor *Sensor) Id() string {
	return sensor.id
}

func (sensor *Sensor) Marshal() (interface{}, error) {
	helper := serialization.NewStructMarshallerHelper()

	helper.Add("_id", sensor.id)
	helper.Add("sensorId", sensor.sensorId)
	helper.Add("deviceClass", sensor.deviceClass)
	helper.Add("stateClass", sensor.stateClass)
	helper.Add("unitOfMeasurement", sensor.unitOfMeasurement)
	helper.Add("accuracyDecimals", sensor.accuracyDecimals)

	return helper.Build()
}

func notifySensors(session *sessions.Session, arg struct{}) (uint64, error) {
	viewId := notification.NotifyView[*Sensor](session, sensors)
	return viewId, nil
}

func initSensors() {

	sensors = store.NewContainer[*Sensor]("sensors")

	sensors.Set(&Sensor{
		id:                "test",
		sensorId:          "sensorId",
		deviceClass:       "deviceClass",
		stateClass:        "stateClass",
		unitOfMeasurement: "unitOfMeasurement",
		accuracyDecimals:  42,
	})

}
