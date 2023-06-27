package serialization

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/iancoleman/strcase"
	"golang.org/x/exp/constraints"
)

type Marshaller interface {
	Marshal() (interface{}, error)
}

type Unmarshaller interface {
	Unmarshal(raw interface{}) error
}

type StructMarshallerHelper struct {
	fields map[string]interface{}
	err    error
}

// Opaque marshalled value.
// Useful for value that we have to return back as-this.
type Opaque struct {
	value interface{}
}

func (o *Opaque) Marshal() (interface{}, error) {
	return o.value, nil
}

func (o *Opaque) Unmarshal(raw interface{}) error {
	o.value = raw
	return nil
}

func NewStructMarshallerHelper() *StructMarshallerHelper {
	return &StructMarshallerHelper{
		fields: make(map[string]interface{}),
		err:    nil,
	}
}

func (helper *StructMarshallerHelper) Add(key string, value any) {
	if helper.err != nil {
		return
	}

	helper.fields[key], helper.err = Marshal(value)
}

func (helper *StructMarshallerHelper) Build() (interface{}, error) {
	return helper.fields, helper.err
}

func Marshal(value any) (interface{}, error) {
	return marshalValue(reflect.Indirect(reflect.ValueOf(value)))
}

func Unmarshal(raw interface{}, value any) error {
	return unmarshalValue(raw, reflect.Indirect(reflect.ValueOf(value)))
}

func marshalMerge(value reflect.Value, dest map[string]interface{}) error {
	valueType := reflect.TypeOf(value.Interface())

	if valueType.Kind() != reflect.Struct {
		return fmt.Errorf("Cannot marshal-merge value of type '%s'", valueType.String())
	}

	for fieldIndex := 0; fieldIndex < valueType.NumField(); fieldIndex++ {
		field := valueType.Field(fieldIndex)
		if !field.IsExported() {
			continue
		}

		fieldName := strcase.ToLowerCamel(field.Name)
		fieldValue := value.Field(fieldIndex)

		marshaledValue, err := marshalValue(fieldValue)
		if err != nil {
			return err
		}

		dest[fieldName] = marshaledValue
	}

	return nil
}

func unmarshalUnmerge(raw map[string]interface{}, value reflect.Value) error {
	valueType := reflect.TypeOf(value.Interface())

	if valueType.Kind() != reflect.Struct {
		return fmt.Errorf("Cannot unmarshal-unmerge value of type '%s'", valueType.String())
	}

	for fieldIndex := 0; fieldIndex < valueType.NumField(); fieldIndex++ {
		field := valueType.Field(fieldIndex)
		if !field.IsExported() {
			continue
		}

		fieldName := strcase.ToLowerCamel(field.Name)
		fieldValue := value.Field(fieldIndex)

		rawValue, ok := raw[fieldName]
		if !ok {
			return fmt.Errorf("Cannot unmarshal-unmerge value of type '%s': value not found for field '%s'", valueType.String(), fieldName)
		}

		err := unmarshalValue(rawValue, fieldValue)
		if err != nil {
			return err
		}

		delete(raw, fieldName)
	}

	return nil
}

func getIfImplements[T interface{}](value reflect.Value) T {
	if value.CanAddr() && value.Kind() != reflect.Pointer {
		value = value.Addr()
	}

	iface, ok := value.Interface().(T)

	if ok {
		return iface
	} else {
		var noop T
		return noop
	}
}

func marshalValue(value reflect.Value) (interface{}, error) {
	if !value.IsValid() {
		return nil, nil
	}

	valueType := reflect.TypeOf(value.Interface())

	if valueType == nil {
		return nil, nil
	}

	if findPluginByConcreteType(valueType) != nil {
		// There is a dedicated plugin for that, no need to handle this
		return value.Interface(), nil
	}

	if marshaller := getIfImplements[Marshaller](value); marshaller != nil {
		return marshaller.Marshal()
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

	case reflect.Pointer, reflect.Interface:
		if value.IsNil() {
			return nil, nil
		}
		return marshalValue(value.Elem())

	case reflect.Struct:
		dest := make(map[string]interface{})
		err := marshalMerge(value, dest)
		if err != nil {
			return nil, err
		}
		return dest, nil

	case reflect.Slice:
		dest := make([]interface{}, value.Len())
		for index := 0; index < value.Len(); index++ {
			marshaledValue, err := marshalValue(value.Index(index))
			if err != nil {
				return nil, err
			}
			dest[index] = marshaledValue
		}
		return dest, nil
	}

	return nil, fmt.Errorf("Cannot marshal type '%s'", valueType.String())
}

func unmarshalValue(raw interface{}, value reflect.Value) error {
	if !value.CanSet() {
		return errors.New("Unmarshal: cannot set value")
	}

	valueType := reflect.TypeOf(value.Interface())

	if findPluginByConcreteType(valueType) != nil {
		// There is a dedicated plugin for that, no need to handle this
		rawValue := reflect.ValueOf(raw)
		if !rawValue.Type().AssignableTo(valueType) {
			return fmt.Errorf("Cannot unmarshal value of type '%s' from '%s'", valueType.String(), rawValue.Type().String())
		}

		value.Set(rawValue)
		return nil
	}

	if unmarshaller := getIfImplements[Unmarshaller](value); unmarshaller != nil {
		if value.IsNil() {
			value.Set(reflect.New(valueType.Elem()))
			unmarshaller = getIfImplements[Unmarshaller](value)
		}

		return unmarshaller.Unmarshal(raw)
	}

	switch valueType.Kind() {
	case reflect.String:
		return unmarshalTypedValue[string](raw, value)
	case reflect.Bool:
		return unmarshalTypedValue[bool](raw, value)
	case reflect.Float32:
		return unmarshalNumericValue[float32](raw, value)
	case reflect.Float64:
		return unmarshalNumericValue[float64](raw, value)
	case reflect.Int:
		return unmarshalNumericValue[int](raw, value)
	case reflect.Int8:
		return unmarshalNumericValue[int8](raw, value)
	case reflect.Int16:
		return unmarshalNumericValue[int16](raw, value)
	case reflect.Int32:
		return unmarshalNumericValue[int32](raw, value)
	case reflect.Int64:
		return unmarshalNumericValue[int64](raw, value)
	case reflect.Uint:
		return unmarshalNumericValue[uint](raw, value)
	case reflect.Uint8:
		return unmarshalNumericValue[uint8](raw, value)
	case reflect.Uint16:
		return unmarshalNumericValue[uint16](raw, value)
	case reflect.Uint32:
		return unmarshalNumericValue[uint32](raw, value)
	case reflect.Uint64:
		return unmarshalNumericValue[uint64](raw, value)

	case reflect.Pointer, reflect.Interface:
		if raw == nil {
			// Leave it to nil
			value.SetZero()
			return nil
		}

		newValue := reflect.New(valueType.Elem())
		err := unmarshalValue(raw, newValue.Elem())
		if err != nil {
			return err
		}

		value.Set(newValue)
		return nil

	case reflect.Struct:
		rawMap, ok := raw.(map[string]interface{})
		if !ok {
			return fmt.Errorf("Cannot unmarshal value of type '%s' from '%s'", valueType.String(), reflect.TypeOf(raw).String())
		}

		err := unmarshalUnmerge(rawMap, value)
		if err != nil {
			return err
		}

		if len(rawMap) > 0 {
			return fmt.Errorf("Cannot unmarshal value of type '%s' from '%s': some fields are unmarshaled", valueType.String(), reflect.TypeOf(raw).String())
		}

		return nil

	case reflect.Slice:
		rawSlice, ok := raw.([]interface{})
		if !ok {
			return fmt.Errorf("Cannot unmarshal value of type '%s' from '%s'", valueType.String(), reflect.TypeOf(raw).String())
		}

		sliceLen := len(rawSlice)
		sliceValue := reflect.MakeSlice(valueType, sliceLen, sliceLen)

		for index := 0; index < sliceLen; index++ {
			err := unmarshalValue(rawSlice[index], sliceValue.Index(index))
			if err != nil {
				return err
			}
		}

		value.Set(sliceValue)

		return nil
	}

	return fmt.Errorf("Cannot unmarshal value of type '%s'", valueType.String())
}

func unmarshalTypedValue[T any](raw interface{}, value reflect.Value) error {
	typedValue, ok := raw.(T)
	if !ok {
		return fmt.Errorf("Cannot unmarshal value of type '%s' from '%s'", reflect.TypeOf(value.Interface()).String(), reflect.TypeOf(raw).String())
	}

	value.Set(reflect.ValueOf(typedValue))
	return nil
}

func unmarshalNumericValue[T constraints.Integer | constraints.Float](raw interface{}, value reflect.Value) error {
	floatValue, ok := raw.(float64)
	if !ok {
		return fmt.Errorf("Cannot unmarshal value of type '%s' from '%s'", reflect.TypeOf(value.Interface()).String(), reflect.TypeOf(raw).String())
	}

	value.Set(reflect.ValueOf(T(floatValue)))
	return nil
}
