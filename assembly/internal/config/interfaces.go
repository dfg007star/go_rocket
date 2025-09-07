package config

import (
	"time"

	"github.com/IBM/sarama"
)

type LoggerConfig interface {
	Level() string
	AsJson() bool
	EnableOTLP() bool
	OTLPEndpoint() string
	ServiceName() string
	ServiceEnvironment() string
}

type MetricsConfig interface {
	CollectorEndpoint() string
	CollectorInterval() time.Duration
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
