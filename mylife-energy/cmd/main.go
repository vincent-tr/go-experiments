package main

import (
	"mylife-energy/pkg/log"
	"mylife-energy/pkg/services"
)

var logger = log.CreateLogger("main")

func main() {
	services.RunServices([]string{"collector"})
}
