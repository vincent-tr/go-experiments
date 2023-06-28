package entities

import "mylife-tools-server/services/io/serialization"

type Sensor struct {
	id                string
	deviceClass       string
	stateClass        string
	unitOfMeasurement string
	accuracyDecimals  uint
}

func (sensor *Sensor) Id() string {
	return sensor.id
}

func (sensor *Sensor) DeviceClass() string {
	return sensor.deviceClass
}

func (sensor *Sensor) StateClass() string {
	return sensor.stateClass
}

func (sensor *Sensor) UnitOfMeasurement() string {
	return sensor.unitOfMeasurement
}

func (sensor *Sensor) AccuracyDecimals() uint {
	return sensor.accuracyDecimals
}

func (sensor *Sensor) Marshal() (interface{}, error) {
	helper := serialization.NewStructMarshallerHelper()

	helper.Add("_id", sensor.id)
	helper.Add("deviceClass", sensor.deviceClass)
	helper.Add("stateClass", sensor.stateClass)
	helper.Add("unitOfMeasurement", sensor.unitOfMeasurement)
	helper.Add("accuracyDecimals", sensor.accuracyDecimals)

	return helper.Build()
}

type SensorData struct {
	Id                string
	DeviceClass       string
	StateClass        string
	UnitOfMeasurement string
	AccuracyDecimals  uint
}

func NewSensor(data *SensorData) *Sensor {
	return &Sensor{
		id:                data.Id,
		deviceClass:       data.DeviceClass,
		stateClass:        data.StateClass,
		unitOfMeasurement: data.UnitOfMeasurement,
		accuracyDecimals:  data.AccuracyDecimals,
	}
}

func SensorsEqual(a *Sensor, b *Sensor) bool {
	return a.id == b.id &&
		a.deviceClass == b.deviceClass &&
		a.stateClass == b.stateClass &&
		a.unitOfMeasurement == b.unitOfMeasurement &&
		a.accuracyDecimals == b.accuracyDecimals
}
