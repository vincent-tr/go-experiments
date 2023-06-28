package api

import (
	"mylife-energy/pkg/entities"
	"mylife-energy/pkg/services/live"
	"mylife-tools-server/services/api"
	"mylife-tools-server/services/notification"
	"mylife-tools-server/services/sessions"
)

var liveDef = api.MakeDefinition("live", notifyMeasures, notifyDevices)

func notifyMeasures(session *sessions.Session, arg struct{}) (uint64, error) {
	measures := live.GetMeasures()
	viewId := notification.NotifyView[*entities.Measure](session, measures)
	return viewId, nil
}

func notifyDevices(session *sessions.Session, arg struct{}) (uint64, error) {
	devices := live.GetDevices()
	viewId := notification.NotifyView[*entities.LiveDevice](session, devices)
	return viewId, nil
}
