package io

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type serializedData[T any] struct {
	Type  string `json:"__type"`
	Value T      `json:"value"`
}

// Note: move to metadata?

type Error struct {
	error
}

type errorData struct {
	Message    string `json:"message"`
	StackTrace string `json:"stacktrace"`
}

func (obj Error) MarshalJSON() ([]byte, error) {
	return marshal("error", errorData{
		Message: obj.Error(),
		// TODO: use stacktrace
		StackTrace: obj.Error(),
	})
}

func (obj *Error) UnmarshalJSON(data []byte) error {
	errData := errorData{}

	err := unmarshal(data, "error", &errData)
	if err != nil {
		return err
	}

	// TODO: use stacktrace
	obj.error = errors.New(errData.Message)

	return nil
}

type Time struct {
	time.Time
}

func (obj Time) MarshalJSON() ([]byte, error) {
	msec := obj.UnixMilli()
	return marshal("date", msec)
}

func (obj *Time) UnmarshalJSON(data []byte) error {
	var msec int64

	err := unmarshal(data, "date", &msec)
	if err != nil {
		return err
	}

	*obj = Time{time.UnixMilli(msec)}
	return nil
}

type Buffer []byte

func (obj Buffer) MarshalJSON() ([]byte, error) {
	str := base64.StdEncoding.EncodeToString(obj)
	return marshal("buffer", str)
}

func (obj *Buffer) UnmarshalJSON(data []byte) error {
	var str string

	err := unmarshal(data, "buffer", &str)
	if err != nil {
		return err
	}

	buffer, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return err
	}

	*obj = buffer
	return nil
}

func marshal[T any](dataType string, data T) ([]byte, error) {
	rawObj := serializedData[T]{
		Type:  dataType,
		Value: data,
	}

	return json.Marshal(&rawObj)
}

func unmarshal[T any](raw []byte, dataType string, data *T) error {
	rawObj := serializedData[*T]{Value: data}

	err := json.Unmarshal(raw, &rawObj)
	if err != nil {
		return err
	}

	if rawObj.Type != dataType {
		return errors.New(fmt.Sprintf("Bad object type '%s' (expected '%s')", rawObj.Type, dataType))
	}

	return nil
}
