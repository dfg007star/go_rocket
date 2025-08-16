package env

import (
	"github.com/caarlos0/env/v11"
)

type orderPaidConsumerEnvConfig struct {
	TopicName string `env:"ORDER_PAID_TOPIC_NAME,required"`
	GroupId   string `env:"ORDER_PAID_CONSUMER_GROUP_ID,required"`
}

type orderPaidConsumerConfig struct {
	raw orderPaidConsumerEnvConfig
	consumerConfig
}

func NewOrderPaidConsumerConfig() (*orderPaidConsumerConfig, error) {
	var raw orderPaidConsumerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderPaidConsumerConfig{raw: raw}, nil
}

func (o *orderPaidConsumerConfig) TopicName() string {
	return o.raw.TopicName
}

func (o *orderPaidConsumerConfig) GroupID() string {
	return o.raw.GroupId
}
