package bus

import (
	"mylife-home-common/defines"
	"mylife-home-common/instance_info"
	"mylife-home-common/log"
	"mylife-home-common/tools"
)

var logger = log.CreateLogger("mylife:home:bus")

type Transport struct {
	client   *client
	metadata *Metadata
}

func NewTransport() *Transport {
	client := newClient(defines.InstanceName())
	transport := &Transport{
		client:   client,
		metadata: newMetadata(client),
	}

	transport.client.OnOnlineChanged().Register(func(online bool) {
		if online {
			transport.publishInstanceInfo()
		}
	})

	instance_info.OnUpdate().Register(func(_ *instance_info.InstanceInfo) {
		if transport.Online() {
			transport.publishInstanceInfo()
		}
	})

	return transport
}

func (transport *Transport) Metadata() *Metadata {
	return transport.metadata
}

func (transport *Transport) OnOnlineChanged(callback *OnlineChangedHandler) tools.CallbackRegistration[bool] {
	return transport.client.OnOnlineChanged()
}

func (transport *Transport) Online() bool {
	return transport.client.Online()
}

func (transport *Transport) publishInstanceInfo() {
	data := instance_info.Get()
	transport.metadata.Set("instance-info", data)
}
