package entities

import (
	"mylife-tools-server/services/io/serialization"
	"sort"

	"github.com/gookit/goutil/errorx/panics"
)

type LiveDevice struct {
	id      string
	display string
	typ     DeviceType
	sensors []LiveSensor
}

func (device *LiveDevice) Id() string {
	return device.id
}

func (device *LiveDevice) Display() string {
	return device.display
}

func (device *LiveDevice) Type() DeviceType {
	return device.typ
}

func (device *LiveDevice) Sensors() []LiveSensor {
	return device.sensors
}

func (device *LiveDevice) Marshal() (interface{}, error) {
	helper := serialization.NewStructMarshallerHelper()

	helper.Add("_id", device.id)
	helper.Add("display", device.display)
	helper.Add("type", device.typ)
	helper.Add("sensors", device.sensors)

	return helper.Build()
}

func LiveDevicesEqual(a *LiveDevice, b *LiveDevice) bool {
	if a.id != b.id ||
		a.display != b.display ||
		a.typ != b.typ {
		return false
	}

	if len(a.sensors) != len(b.sensors) {
		return false
	}

	for i := range a.sensors {
		if !LiveSensorsEqual(&a.sensors[i], &b.sensors[i]) {
			return false
		}
	}

	return true
}

type LiveSensor struct {
	key               string
	display           string
	deviceClass       string
	stateClass        string
	unitOfMeasurement string
	accuracyDecimals  uint
}

func (sensor *LiveSensor) Key() string {
	return sensor.key
}

func (sensor *LiveSensor) Display() string {
	return sensor.display
}

func (sensor *LiveSensor) DeviceClass() string {
	return sensor.deviceClass
}

func (sensor *LiveSensor) StateClass() string {
	return sensor.stateClass
}

func (sensor *LiveSensor) UnitOfMeasurement() string {
	return sensor.unitOfMeasurement
}

func (sensor *LiveSensor) AccuracyDecimals() uint {
	return sensor.accuracyDecimals
}

func (sensor *LiveSensor) Marshal() (interface{}, error) {
	helper := serialization.NewStructMarshallerHelper()

	helper.Add("key", sensor.key)
	helper.Add("display", sensor.display)
	helper.Add("deviceClass", sensor.deviceClass)
	helper.Add("stateClass", sensor.stateClass)
	helper.Add("unitOfMeasurement", sensor.unitOfMeasurement)
	helper.Add("accuracyDecimals", sensor.accuracyDecimals)

	return helper.Build()
}

func LiveSensorsEqual(a *LiveSensor, b *LiveSensor) bool {
	return a.key == b.key &&
		a.display == b.display &&
		a.deviceClass == b.deviceClass &&
		a.stateClass == b.stateClass &&
		a.unitOfMeasurement == b.unitOfMeasurement &&
		a.accuracyDecimals == b.accuracyDecimals
}

type LiveDeviceData struct {
	Id      string
	Display string
	Type    DeviceType
	Sensors []LiveSensorData
}

type LiveSensorData struct {
	Key               string
	Display           string
	DeviceClass       string
	StateClass        string
	UnitOfMeasurement string
	AccuracyDecimals  uint
}

func NewLiveDevice(data *LiveDeviceData) *LiveDevice {
	panics.NotNil(data.Sensors)

	sensors := make([]LiveSensor, 0, len(data.Sensors))
	keys := make(map[string]struct{})

	for _, sensorData := range data.Sensors {
		_, duplicate := keys[sensorData.Key]
		panics.IsFalse(duplicate)

		keys[sensorData.Key] = struct{}{}

		sensors = append(sensors, LiveSensor{
			key:               sensorData.Key,
			display:           sensorData.Display,
			deviceClass:       sensorData.DeviceClass,
			stateClass:        sensorData.StateClass,
			unitOfMeasurement: sensorData.UnitOfMeasurement,
			accuracyDecimals:  sensorData.AccuracyDecimals,
		})
	}

	// Ordering for consistency
	sort.Slice(sensors, func(a, b int) bool {
		return sensors[a].key < sensors[b].key
	})

	return &LiveDevice{
		id:      data.Id,
		display: data.Display,
		typ:     data.Type,
		sensors: sensors,
	}
}
