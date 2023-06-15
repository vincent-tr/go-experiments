package io

import (
	"errors"
	"fmt"
	"reflect"
)

type serializerPlugin interface {
	Type() reflect.Type
	TypeId() string
	Encode(value reflect.Value) interface{}
	Decode(raw interface{}) (reflect.Value, error)
}

type serializationKind uint

const (
	serializationNoop serializationKind = iota
	serializationMap
	serializationSlice
)

var pluginsById = make(map[string]serializerPlugin)
var pluginsByType = make(map[reflect.Type]serializerPlugin)
var nativeTypes = make(map[reflect.Type]serializationKind)

func registerEncoder(plugin serializerPlugin) {
	pluginsById[plugin.TypeId()] = plugin
	pluginsByType[plugin.Type()] = plugin
}

/*
   native handled json types:
   - String
   - Float64
   - Map: map[string]interface{}
   - Slice: []interface{}
   - Bool
   - nil
*/

func init() {
	nativeTypes[getType[string]()] = serializationNoop
	nativeTypes[getType[float64]()] = serializationNoop
	nativeTypes[getType[map[string]interface{}]()] = serializationMap
	nativeTypes[getType[[]interface{}]()] = serializationSlice
	nativeTypes[getType[bool]()] = serializationNoop
}

func serializeValue(value interface{}) interface{} {
	// special case to handle nil
	if value == nil {
		return nil
	}

	valueType := reflect.TypeOf(value)

	plugin, ok := pluginsByType[valueType]
	if ok {
		obj := make(map[string]interface{})
		obj["__type"] = plugin.TypeId()
		obj["value"] = plugin.Encode(reflect.ValueOf(value))
		return obj
	}

	if serKind, ok := nativeTypes[valueType]; ok {
		switch serKind {
		case serializationNoop:
			return value

		case serializationMap:
			mapValue := value.(map[string]interface{})
			obj := make(map[string]interface{})
			for key, value := range mapValue {
				obj[key] = serializeValue(value)
			}
			return obj

		case serializationSlice:
			sliceValue := value.([]interface{})
			slice := make([]interface{}, len(sliceValue))
			for i, value := range sliceValue {
				slice[i] = serializeValue(value)
			}
			return slice
		}
	}

	panic(fmt.Sprintf("Unsupported value found: %+v", value))
}

func deserializeValue(value interface{}) (interface{}, error) {
	// special case to handle nil
	if value == nil {
		return nil, nil
	}

	valueType := reflect.TypeOf(value)

	if serKind, ok := nativeTypes[valueType]; ok {
		switch serKind {
		case serializationNoop:
			return value, nil

		case serializationMap:
			// Test for plugin object
			mapValue := value.(map[string]interface{})
			if pluginType, ok := getPluginType(mapValue); ok {
				plugin, ok := pluginsById[pluginType]
				if !ok {
					return nil, errors.New(fmt.Sprintf("Plugin '%s' not found", plugin))
				}

				pluginValue, ok := mapValue["value"]
				if !ok {
					return nil, errors.New("Plugin without value")
				}

				reflectValue, err := plugin.Decode(pluginValue)
				if err != nil {
					return nil, err
				}

				return reflectValue.Interface(), nil
			}

			obj := make(map[string]interface{})
			for key, value := range mapValue {
				newValue, err := deserializeValue(value)
				if err != nil {
					return nil, err
				}
				obj[key] = newValue
			}
			return obj, nil

		case serializationSlice:
			sliceValue := value.([]interface{})
			slice := make([]interface{}, len(sliceValue))
			for i, value := range sliceValue {
				newValue, err := deserializeValue(value)
				if err != nil {
					return nil, err
				}
				slice[i] = newValue
			}
			return slice, nil
		}
	}

	panic(fmt.Sprintf("Unsupported value found: %+v", value))
}

func getPluginType(mapValue map[string]interface{}) (string, bool) {
	pluginRawType, ok := mapValue["__type"]
	if !ok {
		return "", false
	}

	pluginType, ok := pluginRawType.(string)
	if !ok {
		return "", false
	}

	return pluginType, true
}

func getType[T any]() reflect.Type {
	var ptr *T = nil
	return reflect.TypeOf(ptr).Elem()
}
