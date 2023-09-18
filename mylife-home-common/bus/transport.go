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

func (transport *Transport) RegisterOnOnlineChanged(callback *OnlineChangedHandler) {
	transport.client.RegisterOnOnlineChanged(callback)
}

func (transport *Transport) UnregisterOnOnlineChanged(callback *OnlineChangedHandler) {
	transport.client.UnregisterOnOnlineChanged(callback)
}

func (transport *Transport) Online() bool {
	return transport.client.Online()
}
