package mqtt

import (
	"net/url"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	config "mylife-energy/pkg/config"
	log "mylife-energy/pkg/log"
	serviceRegistry "mylife-energy/pkg/service"
)

var logger = log.CreateLogger("mqtt:client")

func init() {
	serviceRegistry.Register(&MqttService{subscriptions: []*Subscription{}})
}

type BusConfig struct {
	ServerUrl string `mapstructure:"serverUrl"`
}

type Subscription struct {
	topic    string
	callback func(data []byte)
}

type MqttService struct {
	client        mqtt.Client
	subscriptions []*Subscription
}

func (service *MqttService) Init() error {
	busConfig := BusConfig{}
	config.BindStructure("bus", &busConfig)

	// add default port if needed
	serverUrl := busConfig.ServerUrl
	uri, err := url.Parse(serverUrl)
	if err == nil && uri.Port() == "" {
		serverUrl += ":1883"
	}

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
		logger.WithField("error", err).Error("Connection lost")
	}

	opts.AutoReconnect = true
	// opts.ResumeSubs = true # does not work if subscriptions are made before initial connection

	service.client = mqtt.NewClient(opts)

	service.client.Connect()

	return nil
}

func (service *MqttService) Terminate() error {
	service.client.Disconnect(250)
	service.client = nil
	return nil
}

func (service *MqttService) ServiceName() string {
	return "mqtt"
}

func (service *MqttService) Dependencies() []string {
	return []string{}
}

func (service *MqttService) Subscribe(topic string, callback func(data []byte)) {
	subscription := &Subscription{topic, callback}
	service.subscriptions = append(service.subscriptions, subscription)

	if service.client.IsConnected() {
		service.subscribe(subscription)
	}

	logger.WithField("topic", topic).Info("Subscribed to topic")
}

func (service *MqttService) subscribe(subscription *Subscription) {
	service.client.Subscribe(subscription.topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		subscription.callback(msg.Payload())
	})
}

// Shortcuts

func Subscribe(topic string, callback func(data []byte)) {
	getService().Subscribe(topic, callback)
}

func getService() *MqttService {
	return serviceRegistry.GetService[*MqttService]("mqtt")
}
