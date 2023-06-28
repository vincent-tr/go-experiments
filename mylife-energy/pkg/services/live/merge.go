package live

import (
	"fmt"
	"mylife-energy/pkg/entities"
	"mylife-tools-server/services/io"
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

	pendingDeviceUpdate bool
}

func makeMerger(measures store.IContainer[*entities.Measure], sensors store.IContainer[*entities.Sensor]) (*merger, error) {
	devices, err := store.GetCollection[*entities.Device]("devices")
	if err != nil {
		return nil, err
	}

	m := &merger{
		measures:            measures,
		sensors:             sensors,
		devices:             devices,
		liveDevices:         store.NewContainer[*entities.LiveDevice]("live-devices"),
		liveMeasures:        store.NewContainer[*entities.Measure]("live-measures"),
		pendingDeviceUpdate: false,
	}

	m.measuresChangedCallback = func(event *store.Event[*entities.Measure]) {
		m.measuresChanged(event)
	}

	m.sensorsChangedCallback = func(event *store.Event[*entities.Sensor]) {
		m.deviceOrSensorChanged()
	}

	m.devicesChangedCallback = func(event *store.Event[*entities.Device]) {
		m.deviceOrSensorChanged()
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

func (m *merger) deviceOrSensorChanged() {
	if m.pendingDeviceUpdate {
		return
	}

	m.pendingDeviceUpdate = true

	io.SubmitIoTask("live/compute-devices", m.computeDevices)
}

func (m *merger) computeDevices() {
	m.pendingDeviceUpdate = false

	devices := make(map[string]*entities.LiveDeviceData)

	for _, device := range m.devices.List() {
		devices[device.DeviceId()] = &entities.LiveDeviceData{
			Id:      device.DeviceId(),
			Display: device.Display(),
			Type:    device.Type(),
			Parent:  device.Parent(),
			Sensors: make([]entities.LiveSensorData, 0),
		}
	}

	for _, sensor := range m.sensors.List() {
		deviceId, sensorKey, display := sensorData(sensor)
		device, exists := devices[deviceId]

		if !exists {
			continue
		}

		device.Sensors = append(device.Sensors, entities.LiveSensorData{
			Key:               sensorKey,
			Display:           display,
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

func sensorData(sensor *entities.Sensor) (deviceId string, sensorKey string, display string) {
	switch sensor.DeviceClass() {
	case "apparent_power":
		display = "Puissance apparente"
		sensorKey = "apparent-power"

	case "power":
		display = "Puissance réelle"
		sensorKey = "real-power"

	case "current":
		display = "Courant"
		sensorKey = "current"

	case "voltage":
		display = "Tension"
		sensorKey = "voltage"

	default:
		panic(fmt.Sprintf("Unexpected device class '%s' on sensor '%s'", sensor.DeviceClass(), sensor.Id()))
	}

	panics.IsTrue(strings.HasSuffix(sensor.Id(), sensorKey))

	deviceIdLen := len(sensor.Id()) - len(sensorKey) - 1
	deviceId = sensor.Id()[:deviceIdLen]

	return
}
