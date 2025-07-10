package traders

import (
	"go-experiments/brokers"
	"go-experiments/traders/basic"
	"go-experiments/traders/gpt"
)

func SetupBasicTrader(broker brokers.Broker) {
	basic.Setup(broker)
}

func SetupGptTrader(broker brokers.Broker) {
	gpt.Setup(broker)
}
