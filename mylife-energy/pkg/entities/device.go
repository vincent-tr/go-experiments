package entities

import (
	"mylife-tools-server/services/io/serialization"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeviceType string

const (
	Node  DeviceType = "node"
	Group DeviceType = "group"
	Main  DeviceType = "main"
	Solar DeviceType = "solar"
)

type Device struct {
	id       string
	deviceId string
	display  string
	typ      DeviceType
	computed bool
	parent   string
}

func (device *Device) Id() string {
	return device.id
}

func (device *Device) DeviceId() string {
	return device.deviceId
}

func (device *Device) Display() string {
	return device.display
}

func (device *Device) Type() DeviceType {
	return device.typ
}

func (device *Device) Computed() bool {
	return device.computed
}

func (device *Device) Parent() string {
	return device.parent
}

func (device *Device) Marshal() (interface{}, error) {
	helper := serialization.NewStructMarshallerHelper()

	helper.Add("_id", device.id)
	helper.Add("deviceId", device.deviceId)
	helper.Add("display", device.display)
	helper.Add("type", device.typ)
	helper.Add("computed", device.computed)
	helper.Add("parent", device.parent)

	return helper.Build()
}

type deviceData struct {
	Id       primitive.ObjectID
	DeviceId string
	Display  string
	Type     string
	Computed bool
	Parent   string
}

func deviceEncode(device *Device) ([]byte, error) {
	id, err := primitive.ObjectIDFromHex(device.id)
	if err != nil {
		return nil, err
	}

	return bson.Marshal(deviceData{
		Id:       id,
		DeviceId: device.deviceId,
		Display:  device.display,
		Type:     string(device.typ),
		Computed: device.computed,
		Parent:   device.parent,
	})
}

func deviceDecode(raw []byte) (*Device, error) {
	data := deviceData{}
	if err := bson.Unmarshal(raw, &data); err != nil {
		return nil, err
	}

	device := &Device{
		id:       data.Id.Hex(),
		deviceId: data.DeviceId,
		display:  data.Display,
		typ:      DeviceType(data.Type),
		computed: data.Computed,
		parent:   data.Parent,
	}

	return device, nil
}
