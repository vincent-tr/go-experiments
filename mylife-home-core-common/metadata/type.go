package metadata

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Type interface {
	String() string
	Encode(value any) (string, error)
	Decode(raw string) (any, error)
}

var parser = regexp.MustCompile(`([a-z]+)(.*)`)
var rangeParser = regexp.MustCompile(`\[(-?\d+);(-?\d+)\]`)
var enumParser = regexp.MustCompile(`{(.[\w_\-,]+)}`)

func ParseType(value string) (Type, error) {
	matchs := parser.FindStringSubmatch(value)
	if matchs == nil {
		return nil, fmt.Errorf("Invalid type '%s'", value)
	}

	var baseType, args string

	switch len(matchs) {
	case 2:
		baseType = matchs[1]

	case 3:
		baseType = matchs[1]
		args = matchs[2]

	default:
		return nil, fmt.Errorf("Invalid type '%s' (bad match len)", value)
	}

	switch baseType {
	case "range":
		matchs := rangeParser.FindStringSubmatch(args)
		if matchs == nil || len(matchs) != 3 {
			return nil, fmt.Errorf("Invalid type '%s' (bad args)", value)
		}

		min, err := strconv.ParseInt(matchs[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Invalid type '%s' (%f)", value, err)
		}

		max, err := strconv.ParseInt(matchs[2], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Invalid type '%s' (%f)", value, err)
		}

		if min >= max {
			return nil, fmt.Errorf("Invalid type '%s' (min >= mX)", value)
		}

		return &rangeType{min: min, max: max}, nil

	case "text":
		if args != "" {
			return nil, fmt.Errorf("Invalid type '%s' (unexpected args)", value)
		}
		return &textType{}, nil

	case "float":
		if args != "" {
			return nil, fmt.Errorf("Invalid type '%s' (unexpected args)", value)
		}
		return &floatType{}, nil

	case "bool":
		if args != "" {
			return nil, fmt.Errorf("Invalid type '%s' (unexpected args)", value)
		}
		return &boolType{}, nil

	case "enum":
		matchs := enumParser.FindStringSubmatch(args)
		if matchs == nil || len(matchs) != 2 {
			return nil, fmt.Errorf("Invalid type '%s' (bad args)", value)
		}

		values := strings.Split(matchs[1], ",")
		if len(values) < 2 {
			return nil, fmt.Errorf("Invalid type '%s' (bad args)", value)
		}

		return &enumType{values: values}, nil

	case "complex":
		if args != "" {
			return nil, fmt.Errorf("Invalid type '%s' (unexpected args)", value)
		}
		return &complexType{}, nil

	default:
		return nil, fmt.Errorf("Invalid type '%s' (unknown type)", value)
	}
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
