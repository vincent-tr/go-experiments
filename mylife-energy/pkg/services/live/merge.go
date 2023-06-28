package live

import (
	"mylife-energy/pkg/entities"
	"mylife-tools-server/services/store"
	"strings"

	"github.com/gookit/goutil/errorx/panics"
)

type merger struct {
	measures store.IContainer[*entities.Measure]
	sensors  store.IContainer[*entities.Sensor]
	devices  store.IContainer[*entities.Device]

	measuresChangedCallback func(event *store.Event[*entities.Measure])
	sensorsChangedCallback  func(event *store.Event[*entities.Sensor])
	devicesChangedCallback  func(event *store.Event[*entities.Device])

	liveDevices  *store.Container[*entities.LiveDevice]
	liveMeasures *store.Container[*entities.Measure]
}

func makeMerger(measures store.IContainer[*entities.Measure], sensors store.IContainer[*entities.Sensor]) (*merger, error) {
	devices, err := store.GetCollection[*entities.Device]("devices")
	if err != nil {
		return nil, err
	}

	m := &merger{
		measures:     measures,
		sensors:      sensors,
		devices:      devices,
		liveDevices:  store.NewContainer[*entities.LiveDevice]("live-devices"),
		liveMeasures: store.NewContainer[*entities.Measure]("live-measures"),
	}

	m.measuresChangedCallback = func(event *store.Event[*entities.Measure]) {
		m.measuresChanged(event)
	}

	m.sensorsChangedCallback = func(event *store.Event[*entities.Sensor]) {
		m.sensorsChanged(event)
	}

	m.devicesChangedCallback = func(event *store.Event[*entities.Device]) {
		m.devicesChanged(event)
	}

	m.measures.AddListener(&m.measuresChangedCallback)
	m.sensors.AddListener(&m.sensorsChangedCallback)
	m.devices.AddListener(&m.devicesChangedCallback)

	m.computeDevices()

	for _, measure := range m.measures.List() {
		// TODO: compute
		m.liveMeasures.Set(measure)
	}

	return m, nil
}

func (m *merger) terminate() {
	m.measures.RemoveListener(&m.measuresChangedCallback)
	m.sensors.RemoveListener(&m.sensorsChangedCallback)
	m.devices.RemoveListener(&m.devicesChangedCallback)

	m.measures = nil
	m.sensors = nil
	m.devices = nil
	m.liveDevices = nil
	m.liveMeasures = nil
}

func (m *merger) measuresChanged(event *store.Event[*entities.Measure]) {
	// TODO: compute

	switch event.Type() {
	case store.Create, store.Update:
		m.liveMeasures.Set(event.After())
	case store.Remove:
		m.liveMeasures.Delete(event.Before().Id())
	}
}

func (m *merger) sensorsChanged(event *store.Event[*entities.Sensor]) {
	m.computeDevices()
}

func (m *merger) devicesChanged(event *store.Event[*entities.Device]) {
	m.computeDevices()
}

func (m *merger) computeDevices() {
	devices := make(map[string]entities.LiveDeviceData)

	for _, device := range m.devices.List() {
		devices[device.DeviceId()] = entities.LiveDeviceData{
			Id:      device.DeviceId(),
			Display: device.Display(),
			Type:    device.Type(),
			Parent:  device.Parent(),
			Sensors: make([]entities.LiveSensorData, 0),
		}
	}

	for _, sensor := range m.sensors.List() {
		deviceId, sensorKey := splitSensorId(sensor.Id())
		device, exists := devices[deviceId]

		if !exists {
			continue
		}

		device.Sensors = append(device.Sensors, entities.LiveSensorData{
			Key:               sensorKey,
			Display:           sensorDisplay(sensor),
			DeviceClass:       sensor.DeviceClass(),
			StateClass:        sensor.StateClass(),
			UnitOfMeasurement: sensor.UnitOfMeasurement(),
			AccuracyDecimals:  sensor.AccuracyDecimals(),
		})
	}

	list := make([]*entities.LiveDevice, 0, len(devices))

	for _, deviceData := range devices {
		list = append(list, entities.NewLiveDevice(deviceData))
	}

	logger.Debugf("Updating %d devices", len(list))

	m.liveDevices.ReplaceAll(list, entities.LiveDevicesEqual)
}

func splitSensorId(value string) (deviceId string, sensorKey string) {
	// first part is device id, last part sensor key
	index := strings.LastIndex(value, "-")
	panics.IsTrue(index > -1)

	deviceId = value[:index]
	sensorKey = value[index+1:]
	return
}

func sensorDisplay(sensor *entities.Sensor) string {
	switch sensor.DeviceClass() {
	case "current":
		return "Courant"
	// TODO: others
	default:
		return sensor.DeviceClass()
	}
}
