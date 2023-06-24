package api

import (
	"mylife-tools-server/services/sessions"
	"reflect"
)

type Method struct {
	callee reflect.Value
	input  reflect.Type
}

func newMethod(callee Callee) *Method {
	method := &Method{}

	method.callee = reflect.ValueOf(callee)

	calleeType := method.callee.Type()

	sessionPtrType := reflect.TypeOf((*sessions.Session)(nil))
	errorType := reflect.TypeOf((*error)(nil)).Elem()

	if calleeType.Kind() != reflect.Func {
		panic("Callee is not a func")
	} else if calleeType.IsVariadic() {
		panic("Callee is variadic")
	} else if calleeType.NumIn() != 2 {
		panic("Callee has not 2 inputs")
	} else if calleeType.NumOut() != 2 {
		panic("Callee has not 2 outputs")
	} else if calleeType.In(0) != sessionPtrType {
		panic("Callee's first input is not session")
	} else if !calleeType.Out(1).Implements(errorType) {
		panic("Callee's second output is not error")
	}

	method.input = calleeType.In(1)

	return method
}

func (method *Method) InputType() reflect.Type {
	return method.input
}

func (method *Method) Call(session *sessions.Session, input reflect.Value) (any, error) {
	inputValues := []reflect.Value{reflect.ValueOf(session), input}
	outputValues := method.callee.Call(inputValues)
	if len(outputValues) != 2 {
		panic("Invalid outputValues size")
	}

	output := outputValues[0].Interface()
	err := toError(outputValues[1].Interface())

	return output, err
}

func toError(value any) error {
	if value == nil {
		return nil
	} else {
		return value.(error)
	}
}
