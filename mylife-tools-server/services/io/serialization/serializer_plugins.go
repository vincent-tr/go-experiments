package serialization

import (
	"encoding/base64"
	"errors"
	"fmt"
	"reflect"
	"time"
)

func init() {
	registerEncoder(&dateSerializerPlugin{})
	registerEncoder(&errorSerializerPlugin{})
	registerEncoder(&bufferSerializerPlugin{})
}

type dateSerializerPlugin struct {
}

func (plugin *dateSerializerPlugin) Type() reflect.Type {
	return getType[time.Time]()
}

func (plugin *dateSerializerPlugin) TypeId() string {
	return "date"
}

func (plugin *dateSerializerPlugin) Encode(value reflect.Value) (interface{}, error) {
	time := value.Interface().(time.Time)
	msec := time.UnixMilli()
	return float64(msec), nil
}

func (plugin *dateSerializerPlugin) Decode(raw interface{}) (reflect.Value, error) {
	rawValue, ok := raw.(float64)
	if !ok {
		return reflect.Value{}, fmt.Errorf("Bad time value : %+v", raw)
	}

	msec := int64(rawValue)
	time := time.UnixMilli(msec)
	return reflect.ValueOf(time), nil
}

type errorSerializerPlugin struct {
}

func (plugin *errorSerializerPlugin) Type() reflect.Type {
	return getType[error]()
}

func (plugin *errorSerializerPlugin) TypeId() string {
	return "error"
}

func (plugin *errorSerializerPlugin) Encode(value reflect.Value) (interface{}, error) {
	err := value.Interface().(error)
	obj := make(map[string]interface{})

	obj["message"] = err.Error()
	// TODO: use stacktrace
	obj["stacktrace"] = err.Error()

	return obj, nil
}

func (plugin *errorSerializerPlugin) Decode(raw interface{}) (reflect.Value, error) {
	obj, ok := raw.(map[string]interface{})
	if !ok {
		return reflect.Value{}, fmt.Errorf("Bad error value : %+v", raw)
	}

	message := obj["message"].(string)
	// TODO: use stacktrace
	// stacktrace := obj["stacktrace"].(string)

	return reflect.ValueOf(errors.New(message)), nil
}

type bufferSerializerPlugin struct {
}

func (plugin *bufferSerializerPlugin) Type() reflect.Type {
	return getType[[]byte]()
}

func (plugin *bufferSerializerPlugin) TypeId() string {
	return "buffer"
}

func (plugin *bufferSerializerPlugin) Encode(value reflect.Value) (interface{}, error) {
	buffer := value.Interface().([]byte)
	raw := base64.StdEncoding.EncodeToString(buffer)
	return raw, nil
}

func (plugin *bufferSerializerPlugin) Decode(raw interface{}) (reflect.Value, error) {
	str, ok := raw.(string)
	if !ok {
		return reflect.Value{}, fmt.Errorf("Bad buffer value : %+v", raw)
	}

	buffer, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return reflect.Value{}, err
	}

	return reflect.ValueOf(buffer), nil
}
