package metadata

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func testParseType(t *testing.T, str string) {
	typ, err := ParseType(str)
	assert.Nil(t, err)
	if err != nil {
		return
	}

	assert.Equal(t, typ.String(), str)
}

func testEncodeDecodeType[ValueType any](t *testing.T, typeStr string, value ValueType, raw string) {
	typ, err := ParseType(typeStr)
	assert.Nil(t, err)
	if err != nil {
		return
	}

	actualRaw, err := typ.Encode(value)
	assert.Nil(t, err)
	if err != nil {
		return
	}

	assert.Equal(t, actualRaw, raw)

	actualValue, err := typ.Decode(raw)
	assert.Nil(t, err)
	if err != nil {
		return
	}

	assert.Equal(t, actualValue, value)
}

func TestParseRange(t *testing.T) {
	testParseType(t, "range[-12;42]")
}

func TestEncodeDecodeRange(t *testing.T) {
	testEncodeDecodeType[int64](t, "range[-12;42]", 40, "40")
}

func TestParseText(t *testing.T) {
	testParseType(t, "text")
}

func TestEncodeDecodeText(t *testing.T) {
	testEncodeDecodeType[string](t, "text", "value", "value")
}

func TestParseFloat(t *testing.T) {
	testParseType(t, "float")
}

func TestEncodeDecodeFloat(t *testing.T) {
	testEncodeDecodeType[float64](t, "float", 42.42, "42.42")
}

func TestParseBool(t *testing.T) {
	testParseType(t, "bool")
}

func TestEncodeDecodeBool(t *testing.T) {
	testEncodeDecodeType[bool](t, "bool", true, "true")
}

func TestParseEnum(t *testing.T) {
	testParseType(t, "enum{one,two,three}")
}

func TestEncodeDecodeEnum(t *testing.T) {
	testEncodeDecodeType[string](t, "enum{one,two,three}", "one", "one")
}

func TestParseComplex(t *testing.T) {
	testParseType(t, "complex")
}

func TestEncodeDecodeComplex(t *testing.T) {
	value := make(map[string]any)
	value["key1"] = 42.42
	value["key2"] = "value"

	testEncodeDecodeType[any](t, "complex", value, `{"key1":42.42,"key2":"value"}`)
}
