package io

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mylife-tools-server/log"
	"mylife-tools-server/services"
	"mylife-tools-server/services/api"
	"mylife-tools-server/services/sessions"

	"mylife-tools-server/utils"

	"nhooyr.io/websocket"
)

type IOSession struct {
	session      *sessions.Session
	socket       *websocket.Conn
	worker       utils.Worker
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
	Error Error `json:"error"`
}

func makeSession(session *sessions.Session, socket *websocket.Conn) IOSession {
	ioSession := IOSession{
		session:      session,
		socket:       socket,
		writeChannel: make(chan []byte, 5),
		readChannel:  make(chan []byte, 5),
		errorChannel: make(chan error),
	}

	ioSession.worker = utils.InitWorker(ioSession.workerEntry)

	return ioSession
}

func (ioSession *IOSession) workerEntry(exit chan struct{}) {
	stopRead := ioSession.startReadSocket()
	stopWrite := ioSession.startWriteSocket()

	for {
		select {
		case data := <-ioSession.readChannel:
			logger.WithFields(log.Fields{"sessionId": ioSession.session.Id(), "message": data}).Trace("Socket received data")
			ioSession.dispatch(data)

		case err := <-ioSession.errorChannel:
			status := websocket.CloseStatus(err)

			switch status {
			case -1:
				logger.WithFields(log.Fields{"sessionId": ioSession.session.Id(), "error": err}).Error("Socket error")
				continue

			case websocket.StatusNormalClosure:
				logger.WithField("sessionId", ioSession.session.Id()).Info("Socket closed")

			default:
				logger.WithFields(log.Fields{"sessionId": ioSession.session.Id(), "error": err}).Info("Socket closed with error")
			}

			// Avoid deadlock
			go func() {
				sessionService := services.GetService[sessions.SessionService]("sessions")
				sessionService.CloseSession(ioSession.session)
			}()

		case <-exit:
			stopRead()
			stopWrite()
			return
		}
	}
}

func (ioSession *IOSession) startReadSocket() func() {
	ctx, cancel := context.WithCancelCause(context.Background())

	stop := func() {
		cancel(exitCause{})
	}

	go func() {
		for {
			msgType, data, err := ioSession.socket.Read(ctx)

			if errors.Is(err, exitCause{}) {
				return
			} else if msgType != websocket.MessageText {
				ioSession.errorChannel <- errors.New(fmt.Sprintf("Expected message of type text, got %s", msgType.String()))
				continue
			} else if err != nil {
				ioSession.errorChannel <- err
			} else {
				ioSession.readChannel <- data
			}
		}
	}()

	return stop
}

func (ioSession *IOSession) startWriteSocket() func() {
	ctx, cancel := context.WithCancelCause(context.Background())
	exitChannel := make(chan struct{}, 1)

	stop := func() {
		exitChannel <- struct{}{}
		cancel(exitCause{})
	}

	go func() {
		for {
			select {
			case data := <-ioSession.writeChannel:
				logger.WithFields(log.Fields{"sessionId": ioSession.session.Id(), "message": data}).Trace("Socket send data")
				err := ioSession.socket.Write(ctx, websocket.MessageText, data)

				if errors.Is(err, exitCause{}) {
					return
				} else if err != nil {
					ioSession.errorChannel <- err
				}

			case <-exitChannel:
				return
			}
		}
	}()

	return stop
}

func (ioSession *IOSession) Close() {
	ioSession.worker.Terminate()
	ioSession.socket.Close(websocket.StatusNormalClosure, "")
}

func (ioSession *IOSession) send(payloadParts ...any) {
	// merge parts json into one payload
	var merged = make(map[string]any)

	for _, part := range payloadParts {
		partial, err := json.Marshal(part)
		if err != nil {
			logger.WithFields(log.Fields{"sessionId": ioSession.session.Id(), "error": err}).Error("Serialization error")
			return
		}

		err = json.Unmarshal(partial, &merged)
		if err != nil {
			logger.WithFields(log.Fields{"sessionId": ioSession.session.Id(), "error": err}).Error("Serialization error")
			return
		}
	}

	data, err := json.Marshal(merged)

	if err != nil {
		logger.WithFields(log.Fields{"sessionId": ioSession.session.Id(), "error": err}).Error("Serialization error")
		return
	}

	ioSession.writeChannel <- data
}

func (ioSession *IOSession) dispatch(data []byte) {
	var engine payloadEngine
	err := json.Unmarshal(data, &engine)

	if err != nil {
		logger.WithFields(log.Fields{"sessionId": ioSession.session.Id(), "error": err}).Error("Deserialization error")
		return
	}

	if engine.Engine != "call" {
		logger.WithFields(log.Fields{"sessionId": ioSession.session.Id(), "engine": engine.Engine}).Debug("Got message with unexpected engine, ignored")
		return
	}

	var input payloadCallInput
	err = json.Unmarshal(data, &input)

	api := services.GetService[api.ApiService]("api")

	method, err := api.Lookup(input.Service, input.Method)
	if err != nil {
		logger.WithFields(log.Fields{"sessionId": ioSession.session.Id(), "error": err}).Error("Error on api lookup")
		ioSession.replyError(&input, err)
		return
	}

	output, err := method.Call(ioSession.session, data)
	if err != nil {
		logger.WithFields(log.Fields{"sessionId": ioSession.session.Id(), "error": err}).Error("Error on method call")
		ioSession.replyError(&input, err)
		return
	}

	ioSession.reply(&input, output)
}

func (ioSession *IOSession) Notify(notification any) {
	ioSession.send(
		payloadEngine{Engine: "notify"},
		notification,
	)
}

func (ioSession *IOSession) reply(input *payloadCallInput, output any) {
	ioSession.send(
		payloadEngine{Engine: "call"},
		payloadCallOutput{Transaction: input.Transaction},
		output,
	)
}

func (ioSession *IOSession) replyError(input *payloadCallInput, err error) {
	ioSession.send(
		payloadEngine{Engine: "call"},
		payloadCallOutput{Transaction: input.Transaction},
		payloadCallError{Error: Error{err}},
	)
}
