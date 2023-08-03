package metadata

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type Type interface {
	String() string
	Encode(value any) (string, error)
	Decode(raw string) (any, error)
}

func ParseType(value string) (Type, error) {
	return nil, fmt.Errorf("TODO")
}

type rangeType struct {
	min int64
	max int64
}

func (this *rangeType) String() string {
	return fmt.Sprintf("range[%d;%d]", this.min, this.max)
}

func (this *rangeType) Encode(value any) (string, error) {
	// TODO: check value
	return strconv.FormatInt(value.(int64), 10), nil
}

func (this *rangeType) Decode(raw string) (any, error) {
	// TODO: check value
	return strconv.ParseInt(raw, 10, 64)
}

type textType struct {
}

func (this *textType) String() string {
	return "text"
}

func (this *textType) Encode(value any) (string, error) {
	return value.(string), nil
}

func (this *textType) Decode(raw string) (any, error) {
	return raw, nil
}

type floatType struct {
}

func (this *floatType) String() string {
	return "float"
}

func (this *floatType) Encode(value any) (string, error) {
	return strconv.FormatFloat(value.(float64), 'g', -1, 64), nil
}

func (this *floatType) Decode(raw string) (any, error) {
	return strconv.ParseFloat(raw, 64)
}

type boolType struct {
}

func (this *boolType) String() string {
	return "bool"
}

func (this *boolType) Encode(value any) (string, error) {
	return strconv.FormatBool(value.(bool)), nil
}

func (this *boolType) Decode(raw string) (any, error) {
	return strconv.ParseBool(raw)
}

type enumType struct {
	values []string
}

func (this *enumType) String() string {
	return fmt.Sprintf("enum{%s}", strings.Join(this.values, ","))
}

func (this *enumType) Encode(value any) (string, error) {
	// TODO: check value
	return value.(string), nil
}

func (this *enumType) Decode(raw string) (any, error) {
	// TODO: check value
	return raw, nil
}

type complexType struct {
}

func (this *complexType) String() string {
	return "complex"
}

func (this *complexType) Encode(value any) (string, error) {
	b, err := json.Marshal(value)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (this *complexType) Decode(raw string) (any, error) {
	var value any
	err := json.Unmarshal([]byte(raw), &value)
	return value, err
}
