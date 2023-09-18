package bus

import (
	"mylife-home-common/config"
	"mylife-home-common/log"
	"strings"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var logger = log.CreateLogger("mylife:home:bus")

type busConfig struct {
	ServerUrl string `mapstructure:"serverUrl"`
}

type onlineChangedHandler func(bool)
type messageHandler func(string, []byte)

type client struct {
	instanceName      string
	mqtt              mqtt.Client
	online            bool
	onlineSync        sync.RWMutex
	onlineCallbacks   map[*onlineChangedHandler]struct{}
	messageCallbacks  map[*messageHandler]struct{}
	subscriptions     map[string]struct{}
	subscriptionsSync sync.RWMutex
}

func newClient(instanceName string) *client {
	conf := busConfig{}
	config.BindStructure("bus", &conf)

	// Need it in advance
	client := &client{
		instanceName:     instanceName,
		onlineCallbacks:  make(map[*onlineChangedHandler]struct{}),
		messageCallbacks: make(map[*messageHandler]struct{}),
		subscriptions:    make(map[string]struct{}),
	}

	options := mqtt.NewClientOptions()
	options.AddBroker(conf.ServerUrl)
	options.SetClientID(instanceName)
	options.SetCleanSession(true)
	options.SetResumeSubs(false)
	options.SetBinaryWill(client.BuildTopic("online"), []byte{}, 0, true)

	options.SetConnectionLostHandler(func(_ mqtt.Client, err error) {
		client.onConnectionLost(err)
	})

	options.SetOnConnectHandler(func(_ mqtt.Client) {
		client.onConnect()
	})

	options.SetDefaultPublishHandler(func(_ mqtt.Client, m mqtt.Message) {
		client.onMessage(m)
	})

	client.mqtt = mqtt.NewClient(options)

	// Note: with auto-retry, the connection may not fail, or for a severe reason (eg: bad config)
	fireAndForget(func() error {
		return tokenSync(client.mqtt.Connect())
	})

	return client
}

func (client *client) Terminate() {
	if client.mqtt.IsConnected() {
		if err := client.ClearRetain(client.BuildTopic("online")); err != nil {
			logger.WithError(err).Error("Send offline error")
		}

		if err := client.clearResidentState(); err != nil {
			logger.WithError(err).Error("Clear resident state error")
		}
	}

	client.mqtt.Disconnect(100)
}

func (client *client) onConnectionLost(err error) {
	l := logger
	if err != nil {
		l = l.WithError(err)
	}

	l.Error("connection lost")

	client.onlineChanged(false)
}

func (client *client) onConnect() {
	fireAndForget(func() error {
		// given the spec, it is unclear if LWT should be executed in case of client takeover, so we run it to be sure
		if err := client.ClearRetain(client.BuildTopic("online")); err != nil {
			return err
		}

		if err := client.clearResidentState(); err != nil {
			return err
		}

		if err := client.Publish(client.BuildTopic("online"), encoding.WriteBool(true), true); err != nil {
			return err
		}

		client.onlineChanged(true)

		topics := client.prepareResubscribe()
		if len(topics) > 0 {
			if err := tokenSync(client.mqtt.SubscribeMultiple(topics, nil)); err != nil {
				return err
			}
		}

		return nil
	})
}

func (client *client) prepareResubscribe() map[string]byte {
	m := make(map[string]byte)

	client.subscriptionsSync.RLock()
	defer client.subscriptionsSync.RUnlock()

	for topic := range client.subscriptions {
		m[topic] = 0 // Topic => QoS
	}

	return m
}

func (client *client) RegisterOnMessage(callback *messageHandler) {
	client.messageCallbacks[callback] = struct{}{}
}

func (client *client) UnregisterOnMessage(callback *messageHandler) {
	delete(client.messageCallbacks, callback)
}

func (client *client) onMessage(m mqtt.Message) {
	for callback := range client.messageCallbacks {
		(*callback)(m.Topic(), m.Payload())
	}
}

func (client *client) RegisterOnOnlineChanged(callback *onlineChangedHandler) {
	client.onlineCallbacks[callback] = struct{}{}
}

func (client *client) UnregisterOnOnlineChanged(callback *onlineChangedHandler) {
	delete(client.onlineCallbacks, callback)
}

func (client *client) onlineChanged(value bool) {
	client.onlineSync.Lock()
	defer client.onlineSync.Unlock()

	if value == client.online {
		return
	}

	client.online = value
	logger.Infof("online: %b", value)

	for callback := range client.onlineCallbacks {
		(*callback)(value)
	}
}

func (client *client) Online() bool {
	client.onlineSync.RLock()
	defer client.onlineSync.RUnlock()

	return client.online
}

func (client *client) clearResidentState() error {
	// register on self state, and remove on every message received
	// wait 1 sec after last message receive

	newTopic := make(chan struct{}, 100)

	clearTopic := func(_ mqtt.Client, m mqtt.Message) {
		// only clear real retained messages
		if m.Retained() && len(m.Payload()) > 0 && strings.HasPrefix(m.Topic(), client.instanceName+"/") {
			newTopic <- struct{}{}
			fireAndForget(func() error {
				return client.ClearRetain(m.Topic())
			})
		}
	}

	selfStateTopic := client.BuildTopic("#")

	if err := tokenSync(client.mqtt.Subscribe(selfStateTopic, 0, clearTopic)); err != nil {
		return err
	}

	timeout := false
	for !timeout {
		select {
		case <-newTopic:
			// reset timer on new topic

		case <-time.After(time.Second):
			// timeout, exit
			timeout = true
		}
	}

	if err := tokenSync(client.mqtt.Unsubscribe(selfStateTopic)); err != nil {
		return err
	}

	return nil
}

func (client *client) BuildTopic(domain string, args ...string) string {
	finalArgs := append([]string{client.instanceName, domain}, args...)
	return strings.Join(finalArgs, "/")
}

func (client *client) BuildRemoteTopic(targetInstance string, domain string, args ...string) string {
	finalArgs := append([]string{targetInstance, domain}, args...)
	return strings.Join(finalArgs, "/")
}

func (client *client) ClearRetainAsync(topic string) mqtt.Token {
	return client.mqtt.Publish(topic, 0, true, []byte{})
}

func (client *client) ClearRetain(topic string) error {
	return tokenSync(client.ClearRetainAsync(topic))
}

func (client *client) PublishAsync(topic string, payload []byte, retain bool) mqtt.Token {
	return client.mqtt.Publish(topic, 0, retain, payload)
}

func (client *client) Publish(topic string, payload []byte, retain bool) error {
	return tokenSync(client.PublishAsync(topic, payload, retain))
}

func (client *client) SubscribeAsync(topics ...string) mqtt.Token {
	client.subscriptionsAdd(topics)

	if client.Online() {
		m := make(map[string]byte)

		for _, topic := range topics {
			m[topic] = 0 // Topic => QoS
		}

		return client.mqtt.SubscribeMultiple(m, nil)
	}

	return newDoneToken()
}

func (client *client) Subscribe(topics ...string) error {
	return tokenSync(client.SubscribeAsync(topics...))
}

func (client *client) subscriptionsAdd(topics []string) {
	client.subscriptionsSync.Lock()
	defer client.subscriptionsSync.Unlock()

	for _, topic := range topics {
		client.subscriptions[topic] = struct{}{}
	}
}

func (client *client) UnsubscribeAsync(topics ...string) mqtt.Token {
	client.subscriptionsDel(topics)

	if client.Online() {
		return client.mqtt.Unsubscribe(topics...)
	}

	return newDoneToken()
}

func (client *client) Unsubscribe(topics ...string) error {
	return tokenSync(client.UnsubscribeAsync(topics...))
}

func (client *client) subscriptionsDel(topics []string) {
	client.subscriptionsSync.Lock()
	defer client.subscriptionsSync.Unlock()

	for _, topic := range topics {
		delete(client.subscriptions, topic)
	}
}

func fireAndForget(callback func() error) {
	go func() {
		if err := callback(); err != nil {
			logger.WithError(err).Error("Fire and forget failed")
		}
	}()
}

func tokenSync(token mqtt.Token) error {
	token.Wait()
	return token.Error()
}
