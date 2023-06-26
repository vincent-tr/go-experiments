package live

import "mylife-tools-server/services/io/serialization"

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

func (sensor *Sensor) SensorId() string {
	return sensor.sensorId
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
	helper.Add("sensorId", sensor.sensorId)
	helper.Add("deviceClass", sensor.deviceClass)
	helper.Add("stateClass", sensor.stateClass)
	helper.Add("unitOfMeasurement", sensor.unitOfMeasurement)
	helper.Add("accuracyDecimals", sensor.accuracyDecimals)

	return helper.Build()
}
