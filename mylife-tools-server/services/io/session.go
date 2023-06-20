package io

import (
	"context"
	"errors"
	"fmt"
	"mylife-tools-server/log"
	"mylife-tools-server/services/api"
	"mylife-tools-server/services/io/serialization"
	"mylife-tools-server/services/sessions"
	"reflect"

	"mylife-tools-server/utils"

	"nhooyr.io/websocket"
)

type ioSession struct {
	session      *sessions.Session
	socket       *websocket.Conn
	worker       *utils.Worker
	writeChannel chan []byte
	readChannel  chan []byte
	errorChannel chan error
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
	Transaction string `json:"transaction"`
}

type payloadCallOutput struct {
	Transaction string `json:"transaction"`
}

type payloadCallError struct {
	Error error `json:"error"`
}

func makeSession(session *sessions.Session, socket *websocket.Conn) ioSession {
	ios := ioSession{
		session:      session,
		socket:       socket,
		writeChannel: make(chan []byte, 5),
		readChannel:  make(chan []byte, 5),
		errorChannel: make(chan error),
	}

	ios.worker = utils.NewWorker(ios.workerEntry)

	return ios
}

func (ios *ioSession) workerEntry(exit chan struct{}) {
	stopRead := ios.startReadSocket()
	stopWrite := ios.startWriteSocket()

	for {
		select {
		case data := <-ios.readChannel:
			logger.WithFields(log.Fields{"sessionId": ios.session.Id(), "message": data}).Trace("Socket received data")
			ios.dispatch(data)

		case err := <-ios.errorChannel:
			status := websocket.CloseStatus(err)

			switch status {
			case -1:
				logger.WithFields(log.Fields{"sessionId": ios.session.Id(), "error": err}).Error("Socket error")
				continue

			case websocket.StatusNormalClosure:
				logger.WithField("sessionId", ios.session.Id()).Info("Socket closed")

			default:
				logger.WithFields(log.Fields{"sessionId": ios.session.Id(), "error": err}).Info("Socket closed with error")
			}

			// Avoid deadlock
			go func() {
				sessions.CloseSession(ios.session)
			}()

		case <-exit:
			stopRead()
			stopWrite()
			return
		}
	}
}

func (ios *ioSession) startReadSocket() func() {
	ctx, cancel := context.WithCancelCause(context.Background())

	stop := func() {
		cancel(exitCause{})
	}

	go func() {
		for {
			msgType, data, err := ios.socket.Read(ctx)

			if errors.Is(err, exitCause{}) {
				return
			} else if msgType != websocket.MessageText {
				ios.errorChannel <- errors.New(fmt.Sprintf("Expected message of type text, got %s", msgType.String()))
				continue
			} else if err != nil {
				ios.errorChannel <- err
			} else {
				ios.readChannel <- data
			}
		}
	}()

	return stop
}

func (ios *ioSession) startWriteSocket() func() {
	ctx, cancel := context.WithCancelCause(context.Background())
	exitChannel := make(chan struct{}, 1)

	stop := func() {
		exitChannel <- struct{}{}
		cancel(exitCause{})
	}

	go func() {
		for {
			select {
			case data := <-ios.writeChannel:
				logger.WithFields(log.Fields{"sessionId": ios.session.Id(), "message": data}).Trace("Socket send data")
				err := ios.socket.Write(ctx, websocket.MessageText, data)

				if errors.Is(err, exitCause{}) {
					return
				} else if err != nil {
					ios.errorChannel <- err
				}

			case <-exitChannel:
				return
			}
		}
	}()

	return stop
}

func (ios *ioSession) Close() {
	ios.worker.Terminate()
	ios.socket.Close(websocket.StatusNormalClosure, "")
}

func (ios *ioSession) send(payloadParts ...any) {
	jsonObj := serialization.NewJsonObject()
	// merge parts json into one payload
	for _, part := range payloadParts {
		err := jsonObj.Marshal(part)
		if err != nil {
			logger.WithFields(log.Fields{"sessionId": ios.session.Id(), "error": err}).Error("Marshal error")
			return
		}
	}

	data, err := serialization.SerializeJsonObject(jsonObj)

	if err != nil {
		logger.WithFields(log.Fields{"sessionId": ios.session.Id(), "error": err}).Error("Serialize error")
		return
	}

	ios.writeChannel <- data
}

func (ios *ioSession) dispatch(data []byte) {
	jsonObj, err := serialization.DeserializeJsonObject(data)
	if err != nil {
		logger.WithFields(log.Fields{"sessionId": ios.session.Id(), "error": err}).Error("Deserialize error")
		return
	}

	var engine payloadEngine
	err = jsonObj.Unmarshal(&engine)

	if err != nil {
		logger.WithFields(log.Fields{"sessionId": ios.session.Id(), "error": err}).Error("Unmarshal error")
		return
	}

	if engine.Engine != "call" {
		logger.WithFields(log.Fields{"sessionId": ios.session.Id(), "engine": engine.Engine}).Debug("Got message with unexpected engine, ignored")
		return
	}

	var input payloadCallInput
	err = jsonObj.Unmarshal(&input)

	if err != nil {
		logger.WithFields(log.Fields{"sessionId": ios.session.Id(), "error": err}).Error("Unmarshal error")
		return
	}

	method, err := api.Lookup(input.Service, input.Method)
	if err != nil {
		logger.WithFields(log.Fields{"sessionId": ios.session.Id(), "error": err}).Error("Error on api lookup")
		ios.replyError(&input, err)
		return
	}

	methodInput := reflect.New(method.InputType())
	if err := jsonObj.Unmarshal(methodInput.Interface()); err != nil {
		logger.WithFields(log.Fields{"sessionId": ios.session.Id(), "error": err}).Error("Unmarshal error")
		ios.replyError(&input, err)
		return
	}

	output, err := method.Call(ios.session, methodInput.Elem())
	if err != nil {
		logger.WithFields(log.Fields{"sessionId": ios.session.Id(), "error": err}).Error("Error on method call")
		ios.replyError(&input, err)
		return
	}

	ios.reply(&input, output)
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
