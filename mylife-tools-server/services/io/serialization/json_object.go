package serialization

import (
	"encoding/json"
	"reflect"
)

type JsonObject struct {
	fields map[string]interface{}
}

func NewJsonObject() *JsonObject {
	return &JsonObject{fields: make(map[string]interface{})}
}

func DeserializeJsonObject(raw []byte) (*JsonObject, error) {
	encoded := make(map[string]interface{})
	err := json.Unmarshal(raw, &encoded)
	if err != nil {
		return nil, err
	}

	fields, err := deserializeValue(encoded)
	if err != nil {
		return nil, err
	}

	obj := &JsonObject{fields: fields.(map[string]interface{})}
	return obj, nil
}

func SerializeJsonObject(obj *JsonObject) ([]byte, error) {
	encoded, err := serializeValue(obj.fields)
	if err != nil {
		return nil, err
	}

	return json.Marshal(encoded)
}

func (obj *JsonObject) Marshal(value any) error {
	return marshalMerge(reflect.ValueOf(value), obj.fields)
}

func (obj *JsonObject) Unmarshal(value any) error {
	return unmarshalUnmerge(obj.fields, reflect.Indirect(reflect.ValueOf(value)))
}
