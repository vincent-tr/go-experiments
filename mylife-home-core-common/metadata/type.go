package metadata

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

type Type interface {
	String() string
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

		return &RangeType{min: min, max: max}, nil

	case "text":
		if args != "" {
			return nil, fmt.Errorf("Invalid type '%s' (unexpected args)", value)
		}
		return &TextType{}, nil

	case "float":
		if args != "" {
			return nil, fmt.Errorf("Invalid type '%s' (unexpected args)", value)
		}
		return &FloatType{}, nil

	case "bool":
		if args != "" {
			return nil, fmt.Errorf("Invalid type '%s' (unexpected args)", value)
		}
		return &BoolType{}, nil

	case "enum":
		matchs := enumParser.FindStringSubmatch(args)
		if matchs == nil || len(matchs) != 2 {
			return nil, fmt.Errorf("Invalid type '%s' (bad args)", value)
		}

		values := strings.Split(matchs[1], ",")
		if len(values) < 2 {
			return nil, fmt.Errorf("Invalid type '%s' (bad args)", value)
		}

		return &EnumType{values: values}, nil

	case "complex":
		if args != "" {
			return nil, fmt.Errorf("Invalid type '%s' (unexpected args)", value)
		}
		return &ComplexType{}, nil

	default:
		return nil, fmt.Errorf("Invalid type '%s' (unknown type)", value)
	}
}

type RangeType struct {
	min int64
	max int64
}

func (typ *RangeType) String() string {
	return fmt.Sprintf("range[%d;%d]", typ.min, typ.max)
}

func (typ *RangeType) Min() int64 {
	return typ.min
}

func (typ *RangeType) Max() int64 {
	return typ.max
}

func (typ *RangeType) validate(value int64) error {
	if value < typ.min || value > typ.max {
		return fmt.Errorf("Invalid value '%d' for '%s'", value, typ.String())
	}

	return nil
}

type TextType struct {
}

func (typ *TextType) String() string {
	return "text"
}

type FloatType struct {
}

func (typ *FloatType) String() string {
	return "float"
}

type BoolType struct {
}

func (typ *BoolType) String() string {
	return "bool"
}

type EnumType struct {
	values []string
}

func (typ *EnumType) String() string {
	return fmt.Sprintf("enum{%s}", strings.Join(typ.values, ","))
}

func (typ *EnumType) NumValues() int {
	return len(typ.values)
}

func (typ *EnumType) Value(index int) string {
	return typ.values[index]
}

func (typ *EnumType) validate(value string) error {
	if !slices.Contains(typ.values, value) {
		return fmt.Errorf("Invalid value '%s' for '%s'", value, typ.String())
	}

	return nil
}

type ComplexType struct {
}

func (typ *ComplexType) String() string {
	return "complex"
}
