package main

import (
	_ "mylife-energy/pkg/collector"

	"mylife-tools-server/log"
	"mylife-tools-server/services"
)

var logger = log.CreateLogger("mylife:energy:main")

func main() {
	services.RunServices([]string{"collector"})
}
