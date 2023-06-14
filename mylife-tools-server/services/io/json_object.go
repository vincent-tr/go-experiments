package io

import (
	"encoding/json"
)

type jsonObject struct {
	fields map[string]interface{}
}

func newJsonObject() *jsonObject {
	return &jsonObject{fields: make(map[string]interface{})}
}

func deserializeJsonObject(raw []byte) (*jsonObject, error) {
	obj := newJsonObject()

	err := json.Unmarshal(raw, &obj.fields)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func serializeJsonObject(obj *jsonObject) ([]byte, error) {
	return json.Marshal(obj.fields)
}

func (obj *jsonObject) marshal(value any) error {
	partial, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = json.Unmarshal(partial, &obj.fields)
	if err != nil {
		return err
	}

	return nil
}

func (obj *jsonObject) unmarshal(value any) error {
	raw, err := json.Marshal(obj.fields)
	if err != nil {
		return err
	}

	err = json.Unmarshal(raw, value)
	if err != nil {
		return err
	}

	return nil
}

/*
Custom marshal/unmarshal
Avec plugins
Basé sur interface{}
Avec others = map[string]interface{}
Avec bonne casse par défaut
Avec attributes pour overrides
*/
