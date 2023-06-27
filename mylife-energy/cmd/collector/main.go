package main

import (
	_ "mylife-energy/pkg/services/collector"

	"mylife-tools-server/log"
	"mylife-tools-server/services"
)

var logger = log.CreateLogger("mylife:energy:main")

func main() {
	services.RunServices([]string{"collector"})
}

/*

next :
- js init front end pour energy (dans une branche de mylife-apps pour l'instant ?)
- go websocket api (pour lancer en mode dev)
	- comment gerer les interfaces communes
- go web server (pour prod)
- go client packaging (pour prod)
- go moteur de view

*/
