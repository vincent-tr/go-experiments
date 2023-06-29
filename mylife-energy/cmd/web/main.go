package main

import (
	ui "mylife-energy-ui"
	"mylife-energy/pkg/api"
	"mylife-energy/pkg/entities"
	"mylife-energy/pkg/services/tesla"
	"mylife-energy/pkg/services/tesla_wall_connector"
	"mylife-tools-server/log"
	"mylife-tools-server/services"
	_ "mylife-tools-server/services/api"
	_ "mylife-tools-server/services/web"
)

/*
next :
- go web server (pour prod)
- go client packaging (pour prod)
*/

var logger = log.CreateLogger("mylife:energy:test")

func main() {
	args := map[string]interface{}{
		"api":   api.Definitions,
		"store": entities.StoreDef,
		"web":   ui.FS,
	}

	services.RunServices([]string{ /*"test",*/ "web", "live", "store"}, args)
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
