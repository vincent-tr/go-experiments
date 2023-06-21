package mqtt

import (
	"net/url"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"mylife-tools-server/config"
	"mylife-tools-server/log"
	"mylife-tools-server/services"
)

var logger = log.CreateLogger("mylife:server:mqtt")

func init() {
	services.Register(&mqttService{subscriptions: []*subscription{}})
}

type busConfig struct {
	ServerUrl string `mapstructure:"serverUrl"`
}

type subscription struct {
	topic    string
	callback func(topic string, data []byte)
}

type mqttService struct {
	client        mqtt.Client
	subscriptions []*subscription
}

func (service *mqttService) Init(arg interface{}) error {
	busConfig := busConfig{}
	config.BindStructure("bus", &busConfig)

	// add default port if needed
	serverUrl := busConfig.ServerUrl
	uri, err := url.Parse(serverUrl)
	if err == nil && uri.Port() == "" {
		serverUrl += ":1883"
	}

	logger.WithField("serverUrl", serverUrl).Info("Config")

	opts := mqtt.NewClientOptions()
	opts.AddBroker(serverUrl)
	opts.SetClientID("mylife-energy-collector")

	opts.DefaultPublishHandler = func(client mqtt.Client, msg mqtt.Message) {
		logger.WithFields(log.Fields{"message": string(msg.Payload()), "topic": msg.Topic()}).Info("Received unexpected message")
	}

	opts.OnConnect = func(client mqtt.Client) {
		for _, subscription := range service.subscriptions {
			service.subscribe(subscription)
		}

		logger.Info("Connected")
	}

	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		logger.WithError(err).Error("Connection lost")
	}

	opts.AutoReconnect = true
	// opts.ResumeSubs = true # does not work if subscriptions are made before initial connection

	service.client = mqtt.NewClient(opts)

	service.client.Connect()

	return nil
}

func (service *mqttService) Terminate() error {
	service.client.Disconnect(250)
	service.client = nil
	return nil
}

func (service *mqttService) ServiceName() string {
	return "mqtt"
}

func (service *mqttService) Dependencies() []string {
	return []string{}
}

func (service *mqttService) Subscribe(topic string, callback func(topic string, data []byte)) {
	sub := &subscription{topic, callback}
	service.subscriptions = append(service.subscriptions, sub)

	if service.client.IsConnected() {
		service.subscribe(sub)
	}

	logger.WithField("topic", topic).Info("Subscribed to topic")
}

func (service *mqttService) subscribe(sub *subscription) {
	service.client.Subscribe(sub.topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		sub.callback(msg.Topic(), msg.Payload())
	})
}

func getService() *mqttService {
	return services.GetService[*mqttService]("mqtt")
}

// Public access

func Subscribe(topic string, callback func(topic string, data []byte)) {
	getService().Subscribe(topic, callback)
}
