package live

import (
	"fmt"
	"mylife-energy/pkg/entities"
	"mylife-tools-server/services/io"
	"mylife-tools-server/services/store"
	"strings"
	"time"

	"github.com/gookit/goutil/errorx/panics"
	"golang.org/x/exp/slices"
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

	for _, deviceData := range devices {
		device := m.findDeviceById(deviceData.Id)
		panics.NotNil(device)

		if !device.Computed() {
			continue
		}

		var reference *entities.Device

		switch device.Type() {
		case entities.Node:
			// Add same sensors than parent on node computed device
			reference = m.findDeviceById(device.Parent())

		case entities.Solar:
			// Add same sensors than main on solar computed device
			reference = m.findFirstDeviceByType(entities.Main)

		case entities.Total:
			// Add same sensors than any group on total computed device
			reference = m.findFirstDeviceByType(entities.Group)

		default:
			logger.WithField("type", device.Type()).Warn("Unexpected computed device type")
			continue
		}

		panics.NotNil(reference)
		referenceData := devices[reference.DeviceId()]
		panics.NotNil(referenceData)

		deviceData.Sensors = referenceData.Sensors
	}

	// Remove groups (could not before to use them as ref data)
	// there are only for computation
	for _, device := range m.devices.List() {
		if device.Type() == entities.Group && !device.Computed() {
			delete(devices, device.DeviceId())
		}
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
		device := m.findDeviceById(liveDevice.Id())
		if device == nil {
			logger.WithField("id", liveDevice.Id()).Warn("Unmatched device")
			continue
		}

		if device.Computed() {
			switch device.Type() {
			case entities.Node:
				m.computeNodeMeasures(&newMeasures, liveDevice, device)

			case entities.Solar:
				m.computeSolarMeasures(&newMeasures, liveDevice, device)

			case entities.Total:
				m.computeTotalMeasures(&newMeasures, liveDevice, device)

			default:
				logger.WithField("type", device.Type()).Warn("Unexpected computed device type")
			}

			continue
		}

		for _, liveSensor := range liveDevice.Sensors() {
			measure := m.findMeasureValue(liveDevice.Id(), liveSensor.Key())
			if measure != nil {
				newMeasures = append(newMeasures, measure)
			}
		}
	}

	m.liveMeasures.ReplaceAll(newMeasures, entities.MeasuresEqual)
}

func (m *merger) computeNodeMeasures(newMeasures *[]*entities.Measure, liveDevice *entities.LiveDevice, device *entities.Device) {
	// diff between parent and siblings
	parent := m.findDeviceById(device.Parent())
	siblings := m.listDeviceByParent(device.Parent())
	selfIndex := slices.IndexFunc(siblings, func(dev *entities.Device) bool { return dev == device })
	siblings = slices.Delete(siblings, selfIndex, selfIndex+1)

	for _, liveSensor := range liveDevice.Sensors() {
		m.computeMeasureFromParentSiblings(newMeasures, liveDevice.Id(), liveSensor.Key(), parent, siblings, false)
	}
}

func (m *merger) findMeasureValue(deviceId string, sensorKey string) *entities.Measure {
	id := fmt.Sprintf("%s-%s", deviceId, sensorKey)

	measure, exists := m.measures.Find(id)
	if !exists {
		logger.WithField("id", id).Warn("Missing measure")
		return nil
	}

	return measure
}

func (m *merger) computeSolarMeasures(newMeasures *[]*entities.Measure, liveDevice *entities.LiveDevice, device *entities.Device) {
	// diff between main and groups
	parent := m.findFirstDeviceByType(entities.Main)
	siblings := m.listDeviceByType(entities.Group)

	for _, liveSensor := range liveDevice.Sensors() {
		m.computeMeasureFromParentSiblings(newMeasures, liveDevice.Id(), liveSensor.Key(), parent, siblings, true)
	}
}

func (m *merger) computeTotalMeasures(newMeasures *[]*entities.Measure, liveDevice *entities.LiveDevice, device *entities.Device) {
	// sum all groups
	groups := m.listDeviceByType(entities.Group)

	for _, liveSensor := range liveDevice.Sensors() {
		m.computeMeasureFromParentSiblings(newMeasures, liveDevice.Id(), liveSensor.Key(), nil, groups, true)
	}
}

func (m *merger) computeMeasureFromParentSiblings(newMeasures *[]*entities.Measure, deviceId string, sensorKey string, parent *entities.Device, siblings []*entities.Device, producer bool) {
	id := fmt.Sprintf("%s-%s", deviceId, sensorKey)
	var data *entities.MeasureData

	if parent == nil {
		data = &entities.MeasureData{
			Id:        id,
			Sensor:    id,
			Timestamp: time.Now(), // Hopefully siblings will be below
			Value:     0,
		}
	} else {
		parentMeasure := m.findMeasureValue(parent.DeviceId(), sensorKey)
		if parentMeasure == nil {
			return
		}

		data = &entities.MeasureData{
			Id:        id,
			Sensor:    id,
			Timestamp: parentMeasure.Timestamp(),
			Value:     parentMeasure.Value(),
		}
	}

	for _, sibling := range siblings {
		panics.IsFalse(sibling.Computed())

		siblingMeasure := m.findMeasureValue(sibling.DeviceId(), sensorKey)
		if siblingMeasure == nil {
			return
		}

		if siblingMeasure.Timestamp().Before(data.Timestamp) {
			data.Timestamp = siblingMeasure.Timestamp()
		}

		data.Value -= siblingMeasure.Value()
	}

	if producer {
		data.Value = -data.Value
	}

	*newMeasures = append(*newMeasures, entities.NewMeasure(data))
}

func (m *merger) findDeviceById(deviceId string) *entities.Device {
	filteredDevices := m.devices.Filter(func(obj *entities.Device) bool { return obj.DeviceId() == deviceId })
	if len(filteredDevices) != 1 {
		return nil
	}

	return filteredDevices[0]
}

func (m *merger) findFirstDeviceByType(deviceType entities.DeviceType) *entities.Device {
	filteredDevices := m.listDeviceByType(deviceType)
	if len(filteredDevices) < 1 {
		return nil
	}

	return filteredDevices[0]
}

func (m *merger) listDeviceByParent(deviceId string) []*entities.Device {
	return m.devices.Filter(func(obj *entities.Device) bool { return obj.Parent() == deviceId })
}

func (m *merger) listDeviceByType(deviceType entities.DeviceType) []*entities.Device {
	return m.devices.Filter(func(obj *entities.Device) bool { return obj.Type() == deviceType })
}
