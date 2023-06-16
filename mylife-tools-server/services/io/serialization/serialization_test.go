package serialization

import (
	"testing"

	"github.com/gookit/goutil/testutil/assert"
)

type payload struct {
	privateValue int
	Value1       int
	Value2       string
	Value3       bool
	Value4       subPayload
	Value5       *subPayload
	Value6       *subPayload
	Value7       []string
}

type subPayload struct {
	Value string
}

const JSON_VALUE = `{"value1":42,"value2":"toto","value3":true,"value4":{"value":"titi"},"value5":null,"value6":{"value":"toto"},"value7":["titi","tata","toto"]}`

func createPayload() payload {
	return payload{
		privateValue: 12,
		Value1:       42,
		Value2:       "toto",
		Value3:       true,
		Value4:       subPayload{Value: "titi"},
		Value5:       nil,
		Value6:       &subPayload{Value: "toto"},
		Value7:       []string{"titi", "tata", "toto"},
	}
}

func TestMarshal(t *testing.T) {
	val := createPayload()

	obj := NewJsonObject()
	err := obj.Marshal(val)
	if err != nil {
		t.Fatal(err)
	}

	json, err := SerializeJsonObject(obj)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(json))

	assert.Equal(t, JSON_VALUE, string(json))
}

func TestUnmarshal(t *testing.T) {
	obj, err := DeserializeJsonObject([]byte(JSON_VALUE))
	if err != nil {
		t.Fatal(err)
	}

	val := payload{}
	err = obj.Unmarshal(&val)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", val)

	want := createPayload()
	want.privateValue = 0

	assert.Equal(t, want, val)
}
