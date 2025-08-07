package json

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
	var untypedObject map[string]json.RawMessage

	if err := json.Unmarshal(jsonData, &untypedObject); err != nil {
		return empty, err
	}

	// Should have an object with an unique key, which is the type of the object
	if len(untypedObject) != 1 {
		return empty, fmt.Errorf("expected a single object type, got %d keys", len(untypedObject))
	}

	var objectType string
	for key := range untypedObject {
		objectType = key
	}

	objectArgs := untypedObject[objectType]

	parser, exists := r.parsers[objectType]
	if !exists {
		return empty, fmt.Errorf("no parser registered for object type %s", objectType)
	}

	object, err := parser(objectArgs)
	if err != nil {
		return empty, fmt.Errorf("failed to parse object of type %s: %w", objectType, err)
	}

	return object, nil
}
