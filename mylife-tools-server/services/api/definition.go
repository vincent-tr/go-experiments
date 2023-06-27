package api

import (
	"reflect"
	"runtime"
	"strings"
)

type Callee any

type NoReturn *struct{}

type ServiceDefinition struct {
	Name    string
	Methods map[string]Callee
}

func MakeDefinition(name string, callees ...Callee) ServiceDefinition {
	def := ServiceDefinition{
		Name:    name,
		Methods: map[string]Callee{},
	}

	for _, callee := range callees {
		name := getFunctionName(callee)
		if _, ok := def.Methods[name]; ok {
			panic("Method with name already exists")
		}

		def.Methods[name] = callee
	}

	return def
}

// https://stackoverflow.com/questions/7052693/how-to-get-the-name-of-a-function-in-go
func getFunctionName(callee Callee) string {
	strs := strings.Split((runtime.FuncForPC(reflect.ValueOf(callee).Pointer()).Name()), ".")
	return strs[len(strs)-1]
}
