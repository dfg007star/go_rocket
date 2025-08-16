package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	httpClient "github.com/dfg007star/go_rocket/notification/internal/client/http"
	telegramClient "github.com/dfg007star/go_rocket/notification/internal/client/http/telegram"
	"github.com/dfg007star/go_rocket/notification/internal/config"
	kafkaConverter "github.com/dfg007star/go_rocket/notification/internal/converter/kafka"
	"github.com/dfg007star/go_rocket/notification/internal/converter/kafka/decoder"
	"github.com/dfg007star/go_rocket/notification/internal/service"
	orderAssembledConsumer "github.com/dfg007star/go_rocket/notification/internal/service/consumer/order_assembled_consumer"
	orderPaidConsumer "github.com/dfg007star/go_rocket/notification/internal/service/consumer/order_paid_consumer"
	telegramService "github.com/dfg007star/go_rocket/notification/internal/service/telegram"
	"github.com/dfg007star/go_rocket/platform/pkg/closer"
	wrappedKafka "github.com/dfg007star/go_rocket/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/dfg007star/go_rocket/platform/pkg/kafka/consumer"
	"github.com/dfg007star/go_rocket/platform/pkg/logger"
	kafkaMiddleware "github.com/dfg007star/go_rocket/platform/pkg/middleware/kafka"
	"github.com/go-telegram/bot"
)

type diContainer struct {
	telegramService service.TelegramService
	telegramClient  httpClient.TelegramClient
	telegramBot     *bot.Bot

	orderAssembledConsumerService    service.OrderAssembledConsumerService
	orderAssembledKafkaConsumerGroup sarama.ConsumerGroup
	orderAssembledKafkaConsumer      wrappedKafka.Consumer
	orderAssembledKafkaDecoder       kafkaConverter.OrderAssembledDecoder

	orderPaidConsumerService    service.OrderPaidConsumerService
	orderPaidKafkaConsumerGroup sarama.ConsumerGroup
	orderPaidKafkaConsumer      wrappedKafka.Consumer
	orderPaidKafkaDecoder       kafkaConverter.OrderPaidDecoder
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

// TelegramService Service
func (d *diContainer) TelegramService(ctx context.Context) service.TelegramService {
	if d.telegramService == nil {
		d.telegramService = telegramService.NewService(
			d.TelegramClient(ctx),
		)
	}

	return d.telegramService
}

func (d *diContainer) TelegramClient(ctx context.Context) httpClient.TelegramClient {
	if d.telegramClient == nil {
		d.telegramClient = telegramClient.NewClient(d.TelegramBot(ctx))
	}

	return d.telegramClient
}

func (d *diContainer) TelegramBot(ctx context.Context) *bot.Bot {
	if d.telegramBot == nil {
		b, err := bot.New(config.TelegramBotConfig().Token())
		if err != nil {
			panic(fmt.Sprintf("failed to create telegram bot: %s\n", err.Error()))
		}

		d.telegramBot = b
	}

	return d.telegramBot
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

// OrderPaidConsumerService Consumer
func (d *diContainer) OrderPaidConsumerService() service.OrderPaidConsumerService {
	if d.orderPaidConsumerService == nil {
		d.orderPaidConsumerService = orderPaidConsumer.NewService(d.OrderPaidKafkaConsumer(), d.OrderPaidKafkaDecoder())
	}

	return d.orderPaidConsumerService
}

func (d *diContainer) OrderPaidKafkaDecoder() kafkaConverter.OrderPaidDecoder {
	if d.orderPaidKafkaDecoder == nil {
		d.orderPaidKafkaDecoder = decoder.NewOrderPaidDecoder()
	}

	return d.orderPaidKafkaDecoder
}

func (d *diContainer) OrderPaidKafkaConsumer() wrappedKafka.Consumer {
	if d.orderPaidKafkaConsumer == nil {
		d.orderPaidKafkaConsumer = wrappedKafkaConsumer.NewConsumer(
			d.OrderPaidKafkaConsumerGroup(),
			[]string{
				config.AppConfig().OrderPaidConsumer.TopicName(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}

	return d.orderPaidKafkaConsumer
}

func (d *diContainer) OrderPaidKafkaConsumerGroup() sarama.ConsumerGroup {
	if d.orderPaidKafkaConsumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidConsumer.GroupID(),
			config.AppConfig().OrderPaidConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create consumer group: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return d.orderPaidKafkaConsumerGroup.Close()
		})

		d.orderPaidKafkaConsumerGroup = consumerGroup
	}

	return d.orderPaidKafkaConsumerGroup
}
