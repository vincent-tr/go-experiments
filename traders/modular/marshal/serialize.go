package marshal

import "encoding/json"

type ToJsonSpec interface {
	ToJsonSpec() (string, any)
}

func ToJSON(v ToJsonSpec) json.RawMessage {
	key, arg := v.ToJsonSpec()
	var raw json.RawMessage
	var err error

	if arg == nil {
		raw, err = json.Marshal(key)
	} else {
		obj := map[string]any{key: arg}
		raw, err = json.Marshal(obj)
	}

	if err != nil {
		panic("failed to marshal JSON: " + err.Error())
	}
	return raw
}
