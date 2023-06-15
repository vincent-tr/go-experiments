package serialization

import (
	"errors"
	"fmt"
	"reflect"

	"golang.org/x/exp/constraints"
)

/*
Custom marshal/unmarshal
Avec plugins
Basé sur interface{}
Avec others = map[string]interface{}
Avec bonne casse par défaut
Avec attributes pour overrides
*/

/*
   native handled json types:
   - String
   - Float64
   - Map: map[string]interface{}
   - Slice: []interface{}
   - Bool
   - nil
*/

func marshalValue(value reflect.Value) (interface{}, error) {
	if !value.IsValid() {
		return nil, nil
	}

	valueType := value.Type()

	if _, ok := pluginsByType[valueType]; ok {
		// There is a dedicated plugin for that, no need to handle this
		return value.Interface(), nil
	}

	switch valueType.Kind() {
	case reflect.String:
		return value.Interface(), nil
	case reflect.Bool:
		return value.Interface(), nil
	case reflect.Float32:
		return float64(value.Interface().(float32)), nil
	case reflect.Float64:
		return value.Interface(), nil
	case reflect.Int:
		return float64(value.Interface().(int)), nil
	case reflect.Int8:
		return float64(value.Interface().(int8)), nil
	case reflect.Int16:
		return float64(value.Interface().(int16)), nil
	case reflect.Int32:
		return float64(value.Interface().(int32)), nil
	case reflect.Int64:
		return float64(value.Interface().(int64)), nil
	case reflect.Uint:
		return float64(value.Interface().(uint)), nil
	case reflect.Uint8:
		return float64(value.Interface().(uint8)), nil
	case reflect.Uint16:
		return float64(value.Interface().(uint16)), nil
	case reflect.Uint32:
		return float64(value.Interface().(uint32)), nil
	case reflect.Uint64:
		return float64(value.Interface().(uint64)), nil
	}

	// TODO: slice and map, pointer

	return nil, errors.New(fmt.Sprintf("Cannot marshal type '%s'", valueType.String()))
}

func unmarshalValue(raw interface{}, value reflect.Value) error {
	if !value.CanSet() {
		return errors.New("Unmarshal: cannot set value")
	}

	valueType := value.Type()

	if _, ok := pluginsByType[valueType]; ok {
		// There is a dedicated plugin for that, no need to handle this
		rawValue := reflect.ValueOf(raw)
		if valueType != rawValue.Type() {
			return errors.New(fmt.Sprintf("Cannot unmarshal value of type '%s' from '%s'", valueType.String(), rawValue.Type().String()))
		}

		value.Set(rawValue)
		return nil
	}

	switch valueType.Kind() {
	case reflect.String:
		if err := unmarshalTypedValue[string](raw, value); err != nil {
			return err
		}
	case reflect.Bool:
		if err := unmarshalTypedValue[bool](raw, value); err != nil {
			return err
		}
	case reflect.Float32:
		if err := unmarshalNumericValue[float32](raw, value); err != nil {
			return err
		}
	case reflect.Float64:
		if err := unmarshalNumericValue[float64](raw, value); err != nil {
			return err
		}
	case reflect.Int:
		if err := unmarshalNumericValue[int](raw, value); err != nil {
			return err
		}
	case reflect.Int8:
		if err := unmarshalNumericValue[int8](raw, value); err != nil {
			return err
		}
	case reflect.Int16:
		if err := unmarshalNumericValue[int16](raw, value); err != nil {
			return err
		}
	case reflect.Int32:
		if err := unmarshalNumericValue[int32](raw, value); err != nil {
			return err
		}
	case reflect.Int64:
		if err := unmarshalNumericValue[int64](raw, value); err != nil {
			return err
		}
	case reflect.Uint:
		if err := unmarshalNumericValue[uint](raw, value); err != nil {
			return err
		}
	case reflect.Uint8:
		if err := unmarshalNumericValue[uint8](raw, value); err != nil {
			return err
		}
	case reflect.Uint16:
		if err := unmarshalNumericValue[uint16](raw, value); err != nil {
			return err
		}
	case reflect.Uint32:
		if err := unmarshalNumericValue[uint32](raw, value); err != nil {
			return err
		}
	case reflect.Uint64:
		if err := unmarshalNumericValue[uint64](raw, value); err != nil {
			return err
		}
	}

	// TODO: slice and map

	return errors.New(fmt.Sprintf("Cannot unmarshal value of type '%s'", valueType.String()))
}

func unmarshalTypedValue[T any](raw interface{}, value reflect.Value) error {
	typedValue, ok := raw.(T)
	if !ok {
		return errors.New(fmt.Sprintf("Cannot unmarshal value of type '%s' from '%s'", value.Type().String(), getType[T]().String()))
	}

	value.Set(reflect.ValueOf(typedValue))
	return nil
}

func unmarshalNumericValue[T constraints.Integer | constraints.Float](raw interface{}, value reflect.Value) error {
	floatValue, ok := raw.(float64)
	if !ok {
		return errors.New(fmt.Sprintf("Cannot unmarshal value of type '%s' from '%s'", value.Type().String(), getType[T]().String()))
	}

	value.Set(reflect.ValueOf(T(floatValue)))
	return nil
}
