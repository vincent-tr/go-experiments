package components

import "mylife-home-common/bus"

type busPublisher struct {
}

func newBusPublisher(transport *bus.Transport) *busPublisher {
	return &busPublisher{}
}

func (publisher *busPublisher) Terminate() {

}
