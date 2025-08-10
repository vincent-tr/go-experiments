package gridsearch

import (
	"fmt"
)

type ParameterSpace map[string][]interface{}

type Combo map[string]interface{}

func (space ParameterSpace) GenerateCombinations() []Combo {
	keys := make([]string, 0, len(space))
	for k := range space {
		keys = append(keys, k)
	}

	var helper func(int, map[string]interface{})
	results := []Combo{}

	helper = func(index int, current map[string]interface{}) {
		if index == len(keys) {
			combo := make(map[string]interface{}, len(current))
			for k, v := range current {
				combo[k] = v
			}
			results = append(results, combo)
			return
		}

		key := keys[index]
		for _, value := range space[key] {
			current[key] = value
			helper(index+1, current)
		}
	}

	helper(0, map[string]interface{}{})
	return results
}

func comboVal[T any](c Combo, key string) T {
	val, ok := c[key]
	if !ok {
		panic(fmt.Sprintf("Key %s not found", key))
	}

	typed, ok := val.(T)
	if !ok {
		panic(fmt.Sprintf("Value of key %s: wrong type", key))
	}

	return typed
}

func (c Combo) Float(key string) float64 {
	return comboVal[float64](c, key)
}

func (c Combo) Bool(key string) bool {
	return comboVal[bool](c, key)
}

func (c Combo) Int(key string) int {
	return comboVal[int](c, key)
}

func (c Combo) String(key string) string {
	return comboVal[string](c, key)
}
