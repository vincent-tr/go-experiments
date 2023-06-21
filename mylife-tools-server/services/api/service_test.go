package api

import (
	"errors"
	"mylife-tools-server/services/sessions"
	"testing"
)

type input1 struct {
	Id string
}

type input2 struct {
}

type output1 = string

type output2 struct {
	value1 int
	value2 float64
}

func api1(session *sessions.Session, in input1) (output1, error) {
	if in.Id == "42" {
		return "", errors.New("Error")
	}

	return in.Id, nil
}

func api2(session *sessions.Session, in input2) (output2, error) {
	return output2{value1: in.Toto, value2: 42.42}, nil
}

var def = MakeDefinition("test", api1, api2)

func initService() apiService {
	svc := apiService{}
	svc.Init()

	svc.RegisterService(def)

	return svc
}

func Test1(t *testing.T) {
	svc := initService()
	defer svc.Terminate()

	method, err := svc.Lookup("test", "api1")
	if err != nil {
		t.Fatal(err)
	}

	msg := []byte(`{ "service": "test", "method": "api1", "id": "fourty two" }`)

	out, err := method.Call(&sessions.Session{}, msg)

	if err != nil {
		t.Fatalf("Test1: %v", err)
	}

	if out != "fourty two" {
		t.Errorf("got %q, wanted fourty two", out)
	}
}

func Test2(t *testing.T) {
	svc := initService()
	defer svc.Terminate()

	method, err := svc.Lookup("test", "api1")
	if err != nil {
		t.Fatal(err)
	}

	msg := []byte(`{ "service": "test", "method": "api1", "id": "42" }`)

	out, err := method.Call(&sessions.Session{}, msg)

	if err.Error() != "Error" {
		t.Errorf("got %q, wanted Error", err.Error())
	}

	if out != "" {
		t.Errorf("got %q, wanted <empty>", out)
	}
}

func Test3(t *testing.T) {
	svc := initService()
	defer svc.Terminate()

	method, err := svc.Lookup("test", "api2")
	if err != nil {
		t.Fatal(err)
	}

	msg := []byte(`{ "service": "test", "method": "api1", "toto": 42 }`)

	out, err := method.Call(&sessions.Session{}, msg)

	if err != nil {
		t.Fatalf("Test3: %v", err)
	}

	expected := output2{value1: 42, value2: 42.42}
	if out != expected {
		t.Errorf("got %v, wanted %v", out, expected)
	}
}
