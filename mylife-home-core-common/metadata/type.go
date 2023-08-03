package metadata

import (
	"strconv"
	"encoding/json"
)

type Type interface {
	String() string
	Encode(value any) (string, error)
	Decode(raw string) (any, error)
}

func ParseType(value string) Type {

}

type rangeType struct {
	min: int64
	max: int64
}

func (this *rangeType) String() string {
	return fmt.Sprintf("range[%d;%d]", this.min, this.max)
}

func (this *rangeType) Encode(value any) (string, error) {
	// TODO: check value
	return strconv.FormatInt(value.(int64))
}

func (this *rangeType) Decode(raw string) (any, error) {
	// TODO: check value
	return strconv.ParseInt(raw)
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
	return raw
}

type floatType struct {
}

func (this *floatType) String() string {
	return "float"
}

func (this *floatType) Encode(value any) (string, error) {
	return strconv.FormatFloat(value.(float64))
}

func (this *floatType) Decode(raw string) (any, error) {
	return strconv.ParseFloat(raw)
}

type boolType struct {
}

func (this *boolType) String() string {
	return "bool"
}

func (this *boolType) Encode(value any) (string, error) {
	return strconv.FormatBool(value.(bool))
}

func (this *boolType) Decode(raw string) (any, error) {
	return strconv.ParseBool(raw)
}

type enumType struct {
	values: []string
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
	return raw
}

type complexType struct {
}

func (this *complexType) String() string {
	return "complex"
}

func (this *complexType) Encode(value any) (string, error) {
	return json
	// TODO
}

func (this *complexType) Decode(raw string) (any, error) {
	// TODO
}
