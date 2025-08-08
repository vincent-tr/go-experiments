package marshal

import (
	"encoding/json"
	"fmt"
)

type Registry[ObjectType any] struct {
	parsers map[string]func(json.RawMessage) (ObjectType, error)
}

func NewRegistry[ObjectType any]() *Registry[ObjectType] {
	return &Registry[ObjectType]{
		parsers: make(map[string]func(json.RawMessage) (ObjectType, error)),
	}
}

func (r *Registry[ObjectType]) RegisterParser(name string, parser func(json.RawMessage) (ObjectType, error)) {
	if _, exists := r.parsers[name]; exists {
		panic("JSON parser already registered for " + name)
	}
	r.parsers[name] = parser
}

func (r *Registry[ObjectType]) FromJSON(jsonData []byte) (ObjectType, error) {
	var empty ObjectType

	key, args, err := r.parseKeyArgs(jsonData)
	if err != nil {
		return empty, err
	}

	parser, exists := r.parsers[key]
	if !exists {
		return empty, fmt.Errorf("no parser registered for object type %s", key)
	}

	object, err := parser(args)
	if err != nil {
		return empty, fmt.Errorf("failed to parse object of type %s: %w", key, err)
	}

	return object, nil
}

func (r *Registry[ObjectType]) parseKeyArgs(jsonData []byte) (string, json.RawMessage, error) {
	// special case: when no args, can be direct string instead of an object
	var keyWithoutArg string
	if err := json.Unmarshal(jsonData, &keyWithoutArg); err == nil {
		return keyWithoutArg, nil, nil
	}

	var untypedObject map[string]json.RawMessage

	if err := json.Unmarshal(jsonData, &untypedObject); err != nil {
		return "", nil, err
	}

	if len(untypedObject) != 1 {
		return "", nil, fmt.Errorf("expected a single object type, got %d keys", len(untypedObject))
	}

	for key, args := range untypedObject {
		return key, args, nil
	}

	return "", nil, fmt.Errorf("no valid object found in JSON data")
}
