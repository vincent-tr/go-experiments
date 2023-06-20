package io

import (
	"fmt"
	"mylife-tools-server/log"
	"mylife-tools-server/services/api"
	"mylife-tools-server/services/io/serialization"
	"mylife-tools-server/services/sessions"
	"reflect"

	socketio "github.com/vchitai/go-socket.io/v4"
)

type ioSession struct {
	session *sessions.Session
	socket  socketio.Conn
}

type exitCause struct {
}

func (err exitCause) Error() string {
	return "Exit"
}

type payloadEngine struct {
	Engine string `json:"engine"`
}

type payloadCallInput struct {
	Service     string `json:"service"`
	Method      string `json:"method"`
	Transaction uint64 `json:"transaction"`
}

type payloadCallOutput struct {
	Transaction uint64 `json:"transaction"`
}

type payloadCallError struct {
	Error error `json:"error"`
}

func newIoSession(session *sessions.Session, socket socketio.Conn) *ioSession {
	return &ioSession{
		session: session,
		socket:  socket,
	}
}

func (ios *ioSession) Close() {
	if err := ios.socket.Close(); err != nil {
		logger.WithError(err).WithField("sessionId", ios.session.Id()).Error("Error closing socket")
	}
}

func (ios *ioSession) send(payloadParts ...any) {
	jsonObj := serialization.NewJsonObject()
	// merge parts json into one payload
	for _, part := range payloadParts {
		err := jsonObj.Marshal(part)
		if err != nil {
			logger.WithError(err).WithField("sessionId", ios.session.Id()).Error("Marshal error")
			return
		}
	}

	data := serialization.FromJsonObject(jsonObj)
	ios.socket.Emit("message", data)

	logger.WithFields(log.Fields{"sessionId": ios.session.Id(), "data": data}).Debug("Sent message")
}

func (ios *ioSession) dispatch(data map[string]interface{}) {

	logger.WithFields(log.Fields{"sessionId": ios.session.Id(), "data": data}).Debug("Got message")

	jsonObj := serialization.IntoJsonObject(data)

	var engine payloadEngine
	err := jsonObj.Unmarshal(&engine)

	if err != nil {
		logger.WithError(err).WithField("sessionId", ios.session.Id()).Error("Unmarshal error")
		return
	}

	if engine.Engine != "call" {
		logger.WithFields(log.Fields{"sessionId": ios.session.Id(), "engine": engine.Engine}).Debug("Got message with unexpected engine, ignored")
		return
	}

	var input payloadCallInput
	err = jsonObj.Unmarshal(&input)

	if err != nil {
		logger.WithError(err).WithField("sessionId", ios.session.Id()).Error("Unmarshal error")
		return
	}

	method, err := api.Lookup(input.Service, input.Method)
	if err != nil {
		logger.WithError(err).WithField("sessionId", ios.session.Id()).Error("Error on api lookup")
		ios.replyError(&input, err)
		return
	}

	SubmitIoTask(fmt.Sprintf("call/%s/%s", input.Service, input.Method), func() {
		methodInput := reflect.New(method.InputType())
		if err := jsonObj.Unmarshal(methodInput.Interface()); err != nil {
			logger.WithError(err).WithField("sessionId", ios.session.Id()).Error("Unmarshal error")
			ios.replyError(&input, err)
			return
		}

		output, err := method.Call(ios.session, methodInput.Elem())
		if err != nil {
			logger.WithError(err).WithField("sessionId", ios.session.Id()).Error("Error on method call")
			ios.replyError(&input, err)
			return
		}

		ios.reply(&input, output)
	})
}

func (ios *ioSession) notify(notification any) {
	ios.send(
		payloadEngine{Engine: "notify"},
		notification,
	)
}

func (ios *ioSession) reply(input *payloadCallInput, output any) {
	ios.send(
		payloadEngine{Engine: "call"},
		payloadCallOutput{Transaction: input.Transaction},
		output,
	)
}

func (ios *ioSession) replyError(input *payloadCallInput, err error) {
	ios.send(
		payloadEngine{Engine: "call"},
		payloadCallOutput{Transaction: input.Transaction},
		payloadCallError{Error: err},
	)
}
