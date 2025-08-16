package env

import (
	"github.com/caarlos0/env/v11"
)

type orderAssembledConsumerEnvConfig struct {
	TopicName string `env:"ORDER_ASSEMBLED_TOPIC_NAME,required"`
	GroupId   string `env:"ORDER_ASSEMBLED_CONSUMER_GROUP_ID,required"`
}

type orderAssembledConsumerConfig struct {
	raw orderAssembledConsumerEnvConfig
	consumerConfig
}

func NewOrderAssembledConsumerConfig() (*orderAssembledConsumerConfig, error) {
	var raw orderAssembledConsumerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderAssembledConsumerConfig{raw: raw}, nil
}

func (o *orderAssembledConsumerConfig) TopicName() string {
	return o.raw.TopicName
}

func (o *orderAssembledConsumerConfig) GroupID() string {
	return o.raw.GroupId
}
