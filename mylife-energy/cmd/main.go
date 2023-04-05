package main

import (
	"net/url"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	logger "mylife-energy/pkg/logger"

	config "github.com/gookit/config/v2"
	configYaml "github.com/gookit/config/v2/yaml"
)

var log = logger.CreateLogger("main")

type Config struct {
	Mongo string `mapstructure:"mongo"`
	Bus   struct {
		ServerUrl string `mapstructure:"serverUrl"`
	} `mapstructure:"bus"`
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.WithFields(logger.Fields{"message": string(msg.Payload()), "topic": msg.Topic()}).Info("Received message")
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Info("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.WithField("error", err).Info("Connection lost")
}

func main() {
	config.WithOptions(config.ParseEnv)

	// add driver for support yaml content
	config.AddDriver(configYaml.Driver)

	err := config.LoadFiles("config.yaml")
	if err != nil {
		panic(err)
	}

	conf := Config{}
	err = config.Decode(&conf)
	if err != nil {
		panic(err)
	}

	log.WithField("config", conf).Info("Config")

	// add default port if needed
	serverUrl := conf.Bus.ServerUrl
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

	log.WithField("topic", topic).Info("Subscribed to topic")
}
