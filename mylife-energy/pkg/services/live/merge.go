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

	pendingDeviceUpdate  bool
	pendingMeasureUpdate bool
}

func makeMerger(measures store.IContainer[*entities.Measure], sensors store.IContainer[*entities.Sensor]) (*merger, error) {
	devices, err := store.GetCollection[*entities.Device]("devices")
	if err != nil {
		return nil, err
	}

	m := &merger{
		measures:             measures,
		sensors:              sensors,
		devices:              devices,
		liveDevices:          store.NewContainer[*entities.LiveDevice]("live-devices"),
		liveMeasures:         store.NewContainer[*entities.Measure]("live-measures"),
		pendingDeviceUpdate:  false,
		pendingMeasureUpdate: false,
	}

	m.measuresChangedCallback = func(event *store.Event[*entities.Measure]) {
		m.measuresChanged()
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
	m.computeMeasures()

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

	// Recompute measures after device changes
	m.measuresChanged()
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

func (m *merger) measuresChanged() {
	if m.pendingMeasureUpdate {
		return
	}

	m.pendingMeasureUpdate = true

	io.SubmitIoTask("live/compute-measures", m.computeMeasures)
}

func (m *merger) computeMeasures() {
	m.pendingMeasureUpdate = false

	newMeasures := make([]*entities.Measure, 0)

	for _, liveDevice := range m.liveDevices.List() {
		filteredDevices := m.devices.Filter(func(obj *entities.Device) bool { return obj.DeviceId() == liveDevice.Id() })
		if len(filteredDevices) != 1 {
			logger.WithField("id", liveDevice.Id()).Warn("Unmatched device")
			continue
		}

		device := filteredDevices[0]
		if device.Computed() {
			// TODO: computed
			continue
		}

		for _, liveSensor := range liveDevice.Sensors() {
			id := fmt.Sprintf("%s-%s", liveDevice.Id(), liveSensor.Key())

			measure, exists := m.measures.Find(id)
			if !exists {
				logger.WithField("id", id).Warn("Missing measure")
				continue
			}

			newMeasures = append(newMeasures, measure)
		}
	}

	// TODO: computed

	m.liveMeasures.ReplaceAll(newMeasures, entities.MeasuresEqual)
}
