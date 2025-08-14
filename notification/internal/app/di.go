package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/dfg007star/go_rocket/notification/internal/config"
	kafkaConverter "github.com/dfg007star/go_rocket/notification/internal/converter/kafka"
	"github.com/dfg007star/go_rocket/notification/internal/converter/kafka/decoder"
	"github.com/dfg007star/go_rocket/notification/internal/service"
	orderAssembledConsumer "github.com/dfg007star/go_rocket/notification/internal/service/consumer/order_assembled_consumer"
	"github.com/dfg007star/go_rocket/platform/pkg/closer"
	wrappedKafka "github.com/dfg007star/go_rocket/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/dfg007star/go_rocket/platform/pkg/kafka/consumer"
	"github.com/dfg007star/go_rocket/platform/pkg/logger"
	kafkaMiddleware "github.com/dfg007star/go_rocket/platform/pkg/middleware/kafka"
)

type diContainer struct {
	orderAssembledConsumerService    service.OrderAssembledConsumerService
	orderAssembledKafkaConsumerGroup sarama.ConsumerGroup
	orderAssembledKafkaConsumer      wrappedKafka.Consumer
	orderAssembledKafkaDecoder       kafkaConverter.OrderAssembledDecoder

	orderPaidConsumerService    service.OrderAssembledConsumerService
	OrderPaidKafkaConsumerGroup sarama.ConsumerGroup
	OrderPaidKafkaConsumer      wrappedKafka.Consumer
	orderPaidKafkaDecoder       kafkaConverter.OrderAssembledDecoder
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

// OrderAssembledConsumerService Consumer
func (d *diContainer) OrderAssembledConsumerService() service.OrderAssembledConsumerService {
	if d.orderAssembledConsumerService == nil {
		d.orderAssembledConsumerService = orderAssembledConsumer.NewService(d.OrderAssembledKafkaConsumer(), d.OrderAssembledKafkaDecoder())
	}

	return d.orderAssembledConsumerService
}

func (d *diContainer) OrderAssembledKafkaDecoder() kafkaConverter.OrderAssembledDecoder {
	if d.orderAssembledKafkaDecoder == nil {
		d.orderAssembledKafkaDecoder = decoder.NewOrderAssembledDecoder()
	}

	return d.orderAssembledKafkaDecoder
}

func (d *diContainer) OrderAssembledKafkaConsumer() wrappedKafka.Consumer {
	if d.orderAssembledKafkaConsumer == nil {
		d.orderAssembledKafkaConsumer = wrappedKafkaConsumer.NewConsumer(
			d.OrderAssembledKafkaConsumerGroup(),
			[]string{
				config.AppConfig().OrderAssembledConsumer.TopicName(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}

	return d.orderAssembledKafkaConsumer
}

func (d *diContainer) OrderAssembledKafkaConsumerGroup() sarama.ConsumerGroup {
	if d.orderAssembledKafkaConsumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderAssembledConsumer.GroupID(),
			config.AppConfig().OrderAssembledConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create consumer group: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return d.orderAssembledKafkaConsumerGroup.Close()
		})

		d.orderAssembledKafkaConsumerGroup = consumerGroup
	}

	return d.orderAssembledKafkaConsumerGroup
}
