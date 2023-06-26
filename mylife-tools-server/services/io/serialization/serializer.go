package serialization

import (
	"errors"
	"fmt"
	"mylife-tools-server/log"
	"reflect"
)

type serializerPlugin interface {
	Type() reflect.Type
	TypeId() string
	Encode(value reflect.Value) (interface{}, error)
	Decode(raw interface{}) (reflect.Value, error)
}

var logger = log.CreateLogger("mylife:server:io:serialization")

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
	logger.Debugf("Added plugin '%s' for type '%s'", plugin.TypeId(), plugin.Type().String())
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

func serializeValue(value interface{}) (interface{}, error) {
	// special case to handle nil
	if value == nil {
		return nil, nil
	}

	valueType := reflect.TypeOf(value)

	plugin := findPluginByConcreteType(valueType)
	if plugin != nil {
		obj := make(map[string]interface{})
		obj["__type"] = plugin.TypeId()

		pluginValue, err := plugin.Encode(reflect.ValueOf(value))
		if err != nil {
			return nil, err
		}

		obj["value"] = pluginValue
		return obj, nil
	}

	if serKind, ok := nativeTypes[valueType]; ok {
		switch serKind {
		case serializationNoop:
			return value, nil

		case serializationMap:
			mapValue := value.(map[string]interface{})
			obj := make(map[string]interface{})
			for key, value := range mapValue {
				newValue, err := serializeValue(value)
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
				newValue, err := serializeValue(value)
				if err != nil {
					return nil, err
				}
				slice[i] = newValue
			}
			return slice, nil
		}
	}

	return nil, fmt.Errorf("Unsupported value found: '%+v' of type '%s'", value, valueType.String())
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
					return nil, fmt.Errorf("Plugin '%s' not found", plugin)
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

func findPluginByConcreteType(valueType reflect.Type) serializerPlugin {
	// Find exact match
	plugin, ok := pluginsByType[valueType]
	if ok {
		return plugin
	}

	plugin = lookupPluginByConcreteType(valueType)

	// Add it (even if nil) to cache it
	pluginsByType[valueType] = plugin

	if plugin != nil {
		logger.Tracef("Associated plugin '%s' to type '%s'", plugin.TypeId(), valueType.String())
	} else {
		logger.Tracef("Associated plugin '<nil>' to type '%s'", valueType.String())
	}

	return plugin
}

func lookupPluginByConcreteType(valueType reflect.Type) serializerPlugin {
	for pluginType, plugin := range pluginsByType {
		if valueType.AssignableTo(pluginType) {
			return plugin
		}
	}

	return nil
}
