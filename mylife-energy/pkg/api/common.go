package api

import (
	"mylife-tools-server/services/api"
	"mylife-tools-server/services/notification"
	"mylife-tools-server/services/sessions"
)

var commonDef = api.MakeDefinition("common", unnotify)

func unnotify(session *sessions.Session, arg struct{ ViewId uint64 }) (api.NoReturn, error) {
	notification.UnnotifyView(session, arg.ViewId)
	return nil, nil
}
