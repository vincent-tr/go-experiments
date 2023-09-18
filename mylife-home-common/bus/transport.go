package bus

import "mylife-home-common/defines"

type Transport struct {
	client *client
}

func NewTransport() *Transport {
	return &Transport{
		client: newClient(defines.InstanceName()),
	}
}
