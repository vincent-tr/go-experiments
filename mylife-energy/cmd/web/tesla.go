package main

import (
	"context"
	"errors"
	"fmt"
	"math"
	"mylife-tools-server/config"
	"strconv"
	"strings"
	"time"

	"github.com/bogosj/tesla"
)

// https://pkg.go.dev/github.com/bogosj/tesla#WithBaseURL
// https://github.com/bogosj/tesla/blob/v1.1.0/examples/manage_car.go
// https://tesla-api.timdorr.com/api-basics/vehicles

type TeslaConfig struct {
	TokenPath    string `mapstructure:"tokenPath"`
	VIN          string `mapstructure:"vin"`
	HomeLocation string `mapstructure:"homeLocation"` // 'latitude longitude'
}

func initTesla() {
	teslaConfig := TeslaConfig{}
	config.BindStructure("tesla", &teslaConfig)

	homePos, err := parsePosition(teslaConfig.HomeLocation)
	if err != nil {
		logger.WithError(err).Error("Parse home location")
		return
	}

	logger.Debugf("Home position: %+v", homePos)

	client, err := tesla.NewClient(context.TODO(), tesla.WithTokenFile(teslaConfig.TokenPath))
	if err != nil {
		logger.WithError(err).Error("New client")
		return
	}

	vehicles, err := client.Vehicles()
	if err != nil {
		logger.WithError(err).Error("Get vehicles")
		return
	}

	var vehicle *tesla.Vehicle

	for _, veh := range vehicles {
		logger.Debugf("VIN: %s, Name: %s, ID: %d, API version: %d\n", veh.Vin, veh.DisplayName, veh.ID, veh.APIVersion)

		if veh.Vin == teslaConfig.VIN {
			vehicle = veh
		}
	}

	if vehicle == nil {
		logger.Errorf("Vehicle with VIN '%s' not found", teslaConfig.VIN)
		return
	}

	status, err := vehicle.MobileEnabled()
	if err != nil {
		logger.WithError(err).Error("MobileEnabled access failed")
		return
	}

	if !status {
		logger.Error("mobile disabled")
		return
	}

	go func() {

		data, err := vehicle.Data()
		if err != nil {
			logger.WithError(err).Error("Data")
			return
		}

		if data.Error != "" {
			logger.WithError(errors.New(data.Error + " " + data.ErrorDescription)).Error("Data.Error")
			return
		}

		chargeData := ChargeData{
			Timestamp: data.Response.ChargeState.Timestamp.Time,
			Status:    data.Response.ChargeState.ChargingState,
			AtHome:    isAtHome(homePos, &data.Response.DriveState),
			Charger: Charger{
				MaxCurrent: data.Response.ChargeState.ChargerPilotCurrent,
				Current:    data.Response.ChargeState.ChargerActualCurrent,
				Power:      data.Response.ChargeState.ChargerPower,
				Voltage:    data.Response.ChargeState.ChargerVoltage,
			},
			Battery: Battery{
				Level:       data.Response.ChargeState.BatteryLevel,
				TargetLevel: data.Response.ChargeState.ChargeLimitSoc,
			},
			Charge: Charge{
				MinCurrent:     5,
				MaxCurrent:     data.Response.ChargeState.ChargeCurrentRequestMax,
				RequestCurrent: data.Response.ChargeState.ChargeCurrentRequest,
				Current:        data.Response.ChargeState.ChargeAmps,
			},
		}

		logger.Infof("%+v", chargeData)

		// logger.Info(vehicle.Wakeup())
	}()
}

const maxDistance = 50 // meters

type position struct {
	lat  float64
	long float64
}

func parsePosition(pos string) (position, error) {

	parts := strings.Split(pos, " ")
	if len(parts) != 2 {
		return position{}, errors.New(fmt.Sprintf("Invalid position '%s' (split fails)", pos))
	}

	latitude, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return position{}, fmt.Errorf("Invalid position '%s' (parse lat) : %w", pos, err)
	}

	longitude, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return position{}, fmt.Errorf("Invalid position '%s' (parse long) : %w", pos, err)
	}

	return position{lat: latitude, long: longitude}, nil
}

func isAtHome(homePos position, state *tesla.DriveState) bool {
	if state.Speed > 0 {
		return false
	}

	curPos := position{
		lat:  state.Latitude,
		long: state.Longitude,
	}

	dist := distance(homePos, curPos)

	return dist <= maxDistance
}

type ChargeData struct {
	Timestamp time.Time
	Status    string // Charging, Stopped TODO others
	AtHome    bool
	Charger   Charger
	Battery   Battery
	Charge    Charge
}

type Charger struct {
	MaxCurrent int // Max charger current (A)
	Current    int // Actual charger current (A)
	Power      int // Actual charger power (kWh)
	Voltage    int // Actual charger voltage (V)
}

type Battery struct {
	Level       int // Actual battery level (%)
	TargetLevel int // Target battery level (%)
}

type Charge struct {
	MinCurrent     int // Min possible current request (A)
	MaxCurrent     int // Max possible current request (A)
	RequestCurrent int // Requested current (A)
	Current        int // Actual current (A)
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
