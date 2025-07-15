package traders

import (
	"go-experiments/brokers"
	"go-experiments/traders/basic"
	"go-experiments/traders/gpt"
)

type GptConfig = gpt.Config

func SetupBasicTrader(broker brokers.Broker) {
	basic.Setup(broker)
}

func SetupGptTrader(broker brokers.Broker, config *GptConfig) {
	gpt.Setup(broker, config)
}
