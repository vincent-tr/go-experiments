package main

import (
	"net/url"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	config "mylife-energy/pkg/config"
	log "mylife-energy/pkg/log"
)

var logger = log.CreateLogger("main")

type MongoConfig = string

type BusConfig struct {
	ServerUrl string `mapstructure:"serverUrl"`
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	logger.WithFields(log.Fields{"message": string(msg.Payload()), "topic": msg.Topic()}).Info("Received message")
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	logger.Info("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	logger.WithField("error", err).Info("Connection lost")
}

func main() {
	mongoConfig := config.GetString("mongo")

	busConfig := BusConfig{}
	config.BindStructure("bus", &busConfig)

	logger.WithFields(log.Fields{"mongoConfig": mongoConfig, "busConfig": busConfig}).Info("Config")

	// add default port if needed
	serverUrl := busConfig.ServerUrl
	uri, err := url.Parse(serverUrl)
	if err == nil && uri.Port() == "" {
		serverUrl += ":1883"
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(serverUrl)
	opts.SetClientID("mylife-energy-collector")
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	sub(client)

	time.Sleep(30 * time.Second)

	client.Disconnect(250)
}

func sub(client mqtt.Client) {
	topic := "+/energy"
	token := client.Subscribe(topic, 1, nil)
	token.Wait()

	logger.WithField("topic", topic).Info("Subscribed to topic")
}
