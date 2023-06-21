package main

import (
	"mylife-tools-server/log"
	"mylife-tools-server/services"
	"mylife-tools-server/services/api"
	"mylife-tools-server/services/io/serialization"
	"mylife-tools-server/services/notification"
	"mylife-tools-server/services/sessions"
	"mylife-tools-server/services/store"
	_ "mylife-tools-server/services/web"
	"time"
)

var logger = log.CreateLogger("mylife:energy:test")

var sensors *store.Container[*Sensor]

func main() {
	args := make(map[string]interface{})

	sensors = store.NewContainer[*Sensor]("sensors")

	sensors.Set(&Sensor{
		id:                "test",
		sensorId:          "sensorId",
		deviceClass:       "deviceClass",
		stateClass:        "stateClass",
		unitOfMeasurement: "unitOfMeasurement",
		accuracyDecimals:  42,
	})

	args["api"] = []api.ServiceDefinition{
		api.MakeDefinition("common", notifySensors),
	}

	services.RunServices([]string{"test", "web"}, args)
}

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

func notifySensors(session *sessions.Session, arg struct{}) (uint64, error) {
	viewId := notification.NotifyView[*Sensor](session, sensors)
	return viewId, nil
}

func unnotify(session *sessions.Session, arg struct{ viewId uint64 }) error {
	notification.UnnotifyView(session, arg.viewId)
	return nil
}

type testService struct {
}

func (service *testService) Init(arg interface{}) error {
	return nil
}

func (service *testService) Terminate() error {
	return nil
}

func (service *testService) ServiceName() string {
	return "test"
}

func (service *testService) Dependencies() []string {
	return []string{}
}

func init() {
	services.Register(&testService{})
}
