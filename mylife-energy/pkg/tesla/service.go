package tesla

import (
	"context"
	"fmt"
	"math"
	"mylife-tools-server/config"
	"mylife-tools-server/log"
	"mylife-tools-server/services"

	"github.com/bogosj/tesla"
)

var logger = log.CreateLogger("mylife:energy:tesla")

// https://pkg.go.dev/github.com/bogosj/tesla#WithBaseURL
// https://github.com/bogosj/tesla/blob/v1.1.0/examples/manage_car.go
// https://tesla-api.timdorr.com/api-basics/vehicles

// TODO:
// - Status trouver quand debranché
// - Tester vehicle.Wakeup()

type teslaConfig struct {
	TokenPath    string `mapstructure:"tokenPath"`
	VIN          string `mapstructure:"vin"`
	HomeLocation string `mapstructure:"homeLocation"` // 'latitude longitude'
}

type teslaService struct {
	homePos position
	vehicle *tesla.Vehicle
}

func (service *teslaService) Init(arg interface{}) error {
	conf := teslaConfig{}
	config.BindStructure("tesla", &conf)

	homePos, err := parsePosition(conf.HomeLocation)
	if err != nil {
		return fmt.Errorf("Parse home location: %w", err)
	}

	service.homePos = homePos

	logger.Debugf("Home position: %+v", service.homePos)

	client, err := tesla.NewClient(context.TODO(), tesla.WithTokenFile(conf.TokenPath))
	if err != nil {
		return fmt.Errorf("Cannot make new client: %w", err)
	}

	vehicles, err := client.Vehicles()
	if err != nil {
		return fmt.Errorf("Cannot get vehicles: %w", err)
	}

	for _, veh := range vehicles {
		logger.Debugf("VIN: %s, Name: %s, ID: %d, API version: %d\n", veh.Vin, veh.DisplayName, veh.ID, veh.APIVersion)

		if veh.Vin == conf.VIN {
			service.vehicle = veh
		}
	}

	if service.vehicle == nil {
		return fmt.Errorf("Vehicle with VIN '%s' not found", conf.VIN)
	}

	status, err := service.vehicle.MobileEnabled()
	if err != nil {
		return fmt.Errorf("Cannot access MobileEnabled: %w", err)
	}

	if !status {
		return fmt.Errorf("Mobile access disabled")
	}

	return nil
}

func (service *teslaService) Terminate() error {
	service.homePos = position{}
	service.vehicle = nil

	return nil
}

func (service *teslaService) ServiceName() string {
	return "tesla"
}

func (service *teslaService) Dependencies() []string {
	return []string{}
}

func init() {
	services.Register(&teslaService{})
}

func (service *teslaService) fetchChargeData() (*ChargeData, error) {

	data, err := service.vehicle.Data()
	if err != nil {
		return nil, err
	}

	if data.Error != "" {
		return nil, fmt.Errorf("Got data.error: %s %s", data.Error, data.ErrorDescription)
	}

	return newChargeData(&data.Response.ChargeState, service.isAtHome(&data.Response.DriveState)), nil
}

func (service *teslaService) isAtHome(state *tesla.DriveState) bool {
	const maxDistance = 50 // meters

	if state.Speed > 0 {
		return false
	}

	curPos := position{
		lat:  state.Latitude,
		long: state.Longitude,
	}

	dist := distance(service.homePos, curPos)

	return dist <= maxDistance
}

// https://gist.github.com/hotdang-ca/6c1ee75c48e515aec5bc6db6e3265e49
func distance(pos1 position, pos2 position) float64 {
	radlat1 := float64(math.Pi * pos1.lat / 180)
	radlat2 := float64(math.Pi * pos2.lat / 180)

	theta := float64(pos1.long - pos2.long)
	radtheta := float64(math.Pi * theta / 180)

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)
	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / math.Pi
	dist = dist * 60 * 1.1515 * 1.609344 * 1000 // M->KM then KM->M

	return dist
}

func getService() *teslaService {
	return services.GetService[*teslaService]("tesla")
}

// Public access

func FetchChargeData() (*ChargeData, error) {
	return getService().fetchChargeData()
}
