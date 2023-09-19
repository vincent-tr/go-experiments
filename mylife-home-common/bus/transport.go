package bus

import (
	"mylife-home-common/defines"
	"mylife-home-common/instance_info"
	"mylife-home-common/tools"
)

type Transport struct {
	client *client
}

func NewTransport() *Transport {
	transport := &Transport{
		client: newClient(defines.InstanceName()),
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

func (transport *Transport) OnOnlineChanged(callback *OnlineChangedHandler) tools.CallbackRegistration[bool] {
	return transport.client.OnOnlineChanged()
}

func (transport *Transport) Online() bool {
	return transport.client.Online()
}

func (transport *Transport) publishInstanceInfo() {
	data := instance_info.Get()
	// TODO: metadata
	fireAndForget(func() error {
		return transport.client.Publish(transport.client.BuildTopic("metadata", "instance-info"), encoding.WriteJson(data), true)
	})
}
