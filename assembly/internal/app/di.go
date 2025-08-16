package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/dfg007star/go_rocket/assembly/internal/config"
	kafkaConverter "github.com/dfg007star/go_rocket/assembly/internal/converter/kafka"
	"github.com/dfg007star/go_rocket/assembly/internal/converter/kafka/decoder"
	"github.com/dfg007star/go_rocket/assembly/internal/service"
	orderConsumer "github.com/dfg007star/go_rocket/assembly/internal/service/consumer/order_consumer"
	orderProducer "github.com/dfg007star/go_rocket/assembly/internal/service/producer/order_producer"
	"github.com/dfg007star/go_rocket/platform/pkg/closer"
	wrappedKafka "github.com/dfg007star/go_rocket/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/dfg007star/go_rocket/platform/pkg/kafka/consumer"
	wrappedKafkaProducer "github.com/dfg007star/go_rocket/platform/pkg/kafka/producer"
	"github.com/dfg007star/go_rocket/platform/pkg/logger"
	kafkaMiddleware "github.com/dfg007star/go_rocket/platform/pkg/middleware/kafka"
)

type diContainer struct {
	assemblyProducerService service.OrderProducerService
	assemblyConsumerService service.ConsumerService

	assemblyKafkaConsumerGroup sarama.ConsumerGroup
	assemblyKafkaConsumer      wrappedKafka.Consumer
	assemblyKafkaDecoder       kafkaConverter.OrderPaidDecoder

	assemblyKafkaSyncProducer sarama.SyncProducer
	assemblyKafkaProducer     wrappedKafka.Producer
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

// AssemblyProducerService Producer
func (d *diContainer) AssemblyProducerService() service.OrderProducerService {
	if d.assemblyProducerService == nil {
		d.assemblyProducerService = orderProducer.NewService(d.AssemblyKafkaProducer())
	}

	return d.assemblyProducerService
}

func (d *diContainer) AssemblyKafkaProducer() wrappedKafka.Producer {
	if d.assemblyKafkaProducer == nil {
		d.assemblyKafkaProducer = wrappedKafkaProducer.NewProducer(
			d.AssemblyKafkaSyncProducer(),
			config.AppConfig().OrderAssembledProducer.TopicName(),
			logger.Logger(),
		)
	}

	return d.assemblyKafkaProducer
}

func (d *diContainer) AssemblyKafkaSyncProducer() sarama.SyncProducer {
	if d.assemblyKafkaSyncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderAssembledProducer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create sync producer: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka sync producer", func(ctx context.Context) error {
			return p.Close()
		})

		d.assemblyKafkaSyncProducer = p
	}

	return d.assemblyKafkaSyncProducer
}

// AssemblyConsumerService Consumer
func (d *diContainer) AssemblyConsumerService() service.ConsumerService {
	if d.assemblyConsumerService == nil {
		d.assemblyConsumerService = orderConsumer.NewService(d.AssemblyKafkaConsumer(), d.AssemblyKafkaDecoder(), d.AssemblyProducerService())
	}

	return d.assemblyConsumerService
}

func (d *diContainer) AssemblyKafkaDecoder() kafkaConverter.OrderPaidDecoder {
	if d.assemblyKafkaDecoder == nil {
		d.assemblyKafkaDecoder = decoder.NewOrderPaidDecoder()
	}

	return d.assemblyKafkaDecoder
}

func (d *diContainer) AssemblyKafkaConsumer() wrappedKafka.Consumer {
	if d.assemblyKafkaConsumer == nil {
		d.assemblyKafkaConsumer = wrappedKafkaConsumer.NewConsumer(
			d.AssemblyKafkaConsumerGroup(),
			[]string{
				config.AppConfig().OrderPaidConsumer.TopicName(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}

	return d.assemblyKafkaConsumer
}

func (d *diContainer) AssemblyKafkaConsumerGroup() sarama.ConsumerGroup {
	if d.assemblyKafkaConsumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidConsumer.GroupID(),
			config.AppConfig().OrderPaidConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create consumer group: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return d.assemblyKafkaConsumerGroup.Close()
		})

		d.assemblyKafkaConsumerGroup = consumerGroup
	}

	return d.assemblyKafkaConsumerGroup
}
