package main

import (
	"net/url"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	config "github.com/gookit/config/v2"
	configYaml "github.com/gookit/config/v2/yaml"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Mongo string `mapstructure:"mongo"`
	Bus   struct {
		ServerUrl string `mapstructure:"serverUrl"`
	} `mapstructure:"bus"`
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	logger.WithField("message", msg).Info("Received message")
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	logger.Info("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	logger.WithField("error", err).Info("Connection lost")
}

var logger *log.Entry

func init() {
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logger = log.WithFields(log.Fields{
		"name": "main",
	})
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

	logger.WithField("config", conf).Info("Config")

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

	logger.WithField("topic", topic).Info("Subscribed to topic")
}
