package notification

import (
	"fmt"
	"mylife-tools-server/services/io"
	"mylife-tools-server/services/sessions"
	"mylife-tools-server/services/store"
)

type notifySetPayload struct {
	Type   string
	Object any
}

type notifyUnsetPayload struct {
	Type     string
	ObjectId string
}

type notifyPayload struct {
	View int
	List []any
}

type iviewPublisher interface {
	close()
	publish()
}

type viewPublisher[TEntity store.EntityConstraint] struct {
	session  *sessions.Session
	id       int
	view     store.IContainer[TEntity]
	callback func(event *store.Event[TEntity])
	pendings []any
}

func newViewPublisher[TEntity store.EntityConstraint](session *sessions.Session, id int, view store.IContainer[TEntity]) *viewPublisher[TEntity] {
	publisher := &viewPublisher[TEntity]{
		session:  session,
		id:       id,
		view:     view,
		pendings: make([]any, 0),
	}

	publisher.callback = func(event *store.Event[TEntity]) {
		var payload interface{}

		switch event.Type() {
		case store.Create:
		case store.Update:
			payload = &notifySetPayload{Type: "set", Object: event.After()}

		case store.Remove:
			payload = &notifyUnsetPayload{Type: "unset", ObjectId: (*event.Before()).Id()}

		default:
			logger.WithField("eventType", event.Type()).Error("Unexpected event type")
			return
		}

		first := len(publisher.pendings) == 0

		publisher.pendings = append(publisher.pendings, payload)

		// On first push, submit the task
		if first {
			io.SubmitIoTask(fmt.Sprintf("notify/%i", publisher.id), func() {
				payload := &notifyPayload{View: publisher.id, List: publisher.pendings}
				io.NotifySession(*publisher.session, payload)

				publisher.pendings = publisher.pendings[:0]
			})
		}
	}

	publisher.view.AddListener(publisher.callback)

	for obj := range view.List() {
		payload := &notifySetPayload{Type: "set", Object: obj}
		publisher.pendings = append(publisher.pendings, payload)
	}

	return publisher
}

func (publisher *viewPublisher[TEntity]) close() {
	publisher.view.RemoveListener(publisher.callback)

	// publisher.view.Close()
}