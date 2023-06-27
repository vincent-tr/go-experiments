package main

import (
	"mylife-energy/pkg/services/live"
	"mylife-energy/pkg/services/tesla"
	"mylife-energy/pkg/services/tesla_wall_connector"
	"mylife-tools-server/log"
	"mylife-tools-server/services"
	"mylife-tools-server/services/api"
	"mylife-tools-server/services/notification"
	"mylife-tools-server/services/sessions"
	_ "mylife-tools-server/services/web"
)

/*

next :
- js init front end pour energy (dans une branche de mylife-apps pour l'instant ?)
- go websocket api (pour lancer en mode dev)
	- comment gerer les interfaces communes
- go web server (pour prod)
- go client packaging (pour prod)
- go moteur de view

*/

var logger = log.CreateLogger("mylife:energy:test")

func main() {
	args := make(map[string]interface{})

	args["api"] = []api.ServiceDefinition{
		api.MakeDefinition("common", notifySensors, notifyMeasures),
	}

	services.RunServices([]string{"test", "web", "live"}, args)
}

func notifyMeasures(session *sessions.Session, arg struct{}) (uint64, error) {
	measures := live.GetMeasures()
	viewId := notification.NotifyView[*live.Measure](session, measures)
	return viewId, nil
}

func notifySensors(session *sessions.Session, arg struct{}) (uint64, error) {
	sensors := live.GetSensors()
	viewId := notification.NotifyView[*live.Sensor](session, sensors)
	return viewId, nil
}

func unnotify(session *sessions.Session, arg struct{ viewId uint64 }) error {
	notification.UnnotifyView(session, arg.viewId)
	return nil
}

type testService struct {
}

func (service *testService) Init(arg interface{}) error {

	vitals, err := tesla_wall_connector.FetchVitals()
	if err != nil {
		return err
	}

	lifetime, err := tesla_wall_connector.FetchLifetime()
	if err != nil {
		return err
	}

	version, err := tesla_wall_connector.FetchVersion()
	if err != nil {
		return err
	}

	logger.Infof("Vitals: %+v\n", vitals)
	logger.Infof("Lifetime: %+v\n", lifetime)
	logger.Infof("Version: %+v\n", version)

	chargeData, err := tesla.FetchChargeData()
	if err != nil {
		return err
	}

	logger.Infof("ChargeData: %+v\n", chargeData)

	return nil
}

func (service *testService) Terminate() error {
	return nil
}

func (service *testService) ServiceName() string {
	return "test"
}

func (service *testService) Dependencies() []string {
	return []string{"tesla-wall-connector", "tesla"}
}

func init() {
	services.Register(&testService{})
}
