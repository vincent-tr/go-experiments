package traders

import (
	"go-experiments/brokers"
	"go-experiments/traders/basic"
	"go-experiments/traders/gpt"
	"go-experiments/traders/modular"
)

type GptConfig = gpt.Config

func SetupBasicTrader(broker brokers.Broker) {
	basic.Setup(broker)
}

func SetupGptTrader(broker brokers.Broker, config *GptConfig) {
	gpt.Setup(broker, config)
}

func SetupModularTrader(broker brokers.Broker, builder modular.Builder) error {
	return modular.Setup(broker, builder)
}
