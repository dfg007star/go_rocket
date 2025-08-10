package config

import "github.com/IBM/sarama"

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type KafkaConfig interface {
	Brokers() []string
}

type OrderAssembledProducerConfig interface {
	TopicName() string
	Config() *sarama.Config
}

type OrderPaidConsumerConfig interface {
	TopicName() string
	GroupID() string
	Config() *sarama.Config
}
