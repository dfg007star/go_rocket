package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderAPI "github.com/dfg007star/go_rocket/order/internal/api/order/v1"
	grpcClient "github.com/dfg007star/go_rocket/order/internal/client/grpc"
	inventoryServiceClient "github.com/dfg007star/go_rocket/order/internal/client/grpc/inventory/v1"
	paymentServiceClient "github.com/dfg007star/go_rocket/order/internal/client/grpc/payment/v1"
	"github.com/dfg007star/go_rocket/order/internal/config"
	kafkaConverter "github.com/dfg007star/go_rocket/order/internal/converter/kafka"
	"github.com/dfg007star/go_rocket/order/internal/converter/kafka/decoder"
	"github.com/dfg007star/go_rocket/order/internal/repository"
	orderRepository "github.com/dfg007star/go_rocket/order/internal/repository/order"
	"github.com/dfg007star/go_rocket/order/internal/service"
	orderConsumer "github.com/dfg007star/go_rocket/order/internal/service/consumer/order_consumer"
	orderService "github.com/dfg007star/go_rocket/order/internal/service/order"
	orderProducer "github.com/dfg007star/go_rocket/order/internal/service/producer/order_producer"
	"github.com/dfg007star/go_rocket/platform/pkg/closer"
	wrappedKafka "github.com/dfg007star/go_rocket/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/dfg007star/go_rocket/platform/pkg/kafka/consumer"
	wrappedKafkaProducer "github.com/dfg007star/go_rocket/platform/pkg/kafka/producer"
	"github.com/dfg007star/go_rocket/platform/pkg/logger"
	kafkaMiddleware "github.com/dfg007star/go_rocket/platform/pkg/middleware/kafka"
	orderV1 "github.com/dfg007star/go_rocket/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	orderV1API      *orderV1.Server
	orderService    service.OrderService
	orderRepository repository.OrderRepository

	postgresClient  *pgx.Conn
	paymentClient   grpcClient.PaymentClient
	inventoryClient grpcClient.InventoryClient

	orderProducerService   service.OrderProducerService
	orderKafkaProducer     wrappedKafka.Producer
	orderKafkaSyncProducer sarama.SyncProducer

	orderConsumerService    service.ConsumerService
	orderKafkaConsumer      wrappedKafka.Consumer
	orderKafkaConsumerGroup sarama.ConsumerGroup
	orderKafkaDecoder       kafkaConverter.OrderAssembledDecoder
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) OrderV1API(ctx context.Context) *orderV1.Server {
	if d.orderV1API == nil {
		api := orderAPI.NewApi(d.OrderService(ctx))
		orderServer, err := orderV1.NewServer(api)
		if err != nil {
			panic(fmt.Errorf("failed to create order server: %w", err))
		}

		d.orderV1API = orderServer
	}

	return d.orderV1API
}

func (d *diContainer) OrderService(ctx context.Context) service.OrderService {
	if d.orderService == nil {
		d.orderService = orderService.NewOrderService(d.OrderRepository(ctx), d.InventoryClient(ctx), d.PaymentClient(ctx), d.OrderProducerService())
	}

	return d.orderService
}

func (d *diContainer) OrderRepository(ctx context.Context) repository.OrderRepository {
	if d.orderRepository == nil {
		d.orderRepository = orderRepository.NewRepository(d.PostgresClient(ctx))
	}

	return d.orderRepository
}

func (d *diContainer) PostgresClient(ctx context.Context) *pgx.Conn {
	if d.postgresClient == nil {
		con, err := pgx.Connect(ctx, config.AppConfig().Postgres.URI())
		if err != nil {
			panic(fmt.Errorf("failed to connect to database: %w", err))
		}

		closer.AddNamed("Postgres client", func(ctx context.Context) error {
			return con.Close(ctx)
		})

		err = con.Ping(ctx)
		if err != nil {
			panic(fmt.Errorf("database is unavailable: %w", err))
		}

		d.postgresClient = con
	}

	return d.postgresClient
}

func (d *diContainer) PaymentClient(ctx context.Context) grpcClient.PaymentClient {
	if d.paymentClient == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().PaymentGRPC.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(fmt.Errorf("failed to connect to payment grpc client: %w", err))
		}

		closer.AddNamed("Payment GRPC client", func(ctx context.Context) error {
			return conn.Close()
		})

		client := paymentV1.NewPaymentServiceClient(conn)
		d.paymentClient = paymentServiceClient.NewClient(client)
	}

	return d.paymentClient
}

func (d *diContainer) InventoryClient(ctx context.Context) grpcClient.InventoryClient {
	if d.inventoryClient == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().InventoryGRPC.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(fmt.Errorf("failed to connect to inventory grpc client: %w", err))
		}
		closer.AddNamed("Inventory GRPC client", func(ctx context.Context) error {
			return conn.Close()
		})

		client := inventoryV1.NewInventoryServiceClient(conn)
		d.inventoryClient = inventoryServiceClient.NewClient(client)
	}

	return d.inventoryClient
}

// OrderProducerService Producer
func (d *diContainer) OrderProducerService() service.OrderProducerService {
	if d.orderProducerService == nil {
		d.orderProducerService = orderProducer.NewService(d.OrderKafkaProducer())
	}

	return d.orderProducerService
}

func (d *diContainer) OrderKafkaProducer() wrappedKafka.Producer {
	if d.orderKafkaProducer == nil {
		d.orderKafkaProducer = wrappedKafkaProducer.NewProducer(
			d.OrderKafkaSyncProducer(),
			config.AppConfig().OrderPaidProducer.TopicName(),
			logger.Logger(),
		)
	}

	return d.orderKafkaProducer
}

func (d *diContainer) OrderKafkaSyncProducer() sarama.SyncProducer {
	if d.orderKafkaSyncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidProducer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create sync producer: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka sync producer", func(ctx context.Context) error {
			return p.Close()
		})

		d.orderKafkaSyncProducer = p
	}

	return d.orderKafkaSyncProducer
}

// OrderConsumerService Consumer
func (d *diContainer) OrderConsumerService(ctx context.Context) service.ConsumerService {
	if d.orderConsumerService == nil {
		d.orderConsumerService = orderConsumer.NewService(d.OrderKafkaConsumer(), d.OrderKafkaDecoder(), d.OrderRepository(ctx))
	}

	return d.orderConsumerService
}

func (d *diContainer) OrderKafkaDecoder() kafkaConverter.OrderAssembledDecoder {
	if d.orderKafkaDecoder == nil {
		d.orderKafkaDecoder = decoder.NewOrderAssembledDecoder()
	}

	return d.orderKafkaDecoder
}

func (d *diContainer) OrderKafkaConsumer() wrappedKafka.Consumer {
	if d.orderKafkaConsumer == nil {
		d.orderKafkaConsumer = wrappedKafkaConsumer.NewConsumer(
			d.OrderKafkaConsumerGroup(),
			[]string{
				config.AppConfig().OrderAssembledConsumer.TopicName(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}

	return d.orderKafkaConsumer
}

func (d *diContainer) OrderKafkaConsumerGroup() sarama.ConsumerGroup {
	if d.orderKafkaConsumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderAssembledConsumer.GroupID(),
			config.AppConfig().OrderAssembledConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create consumer group: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return d.orderKafkaConsumerGroup.Close()
		})

		d.orderKafkaConsumerGroup = consumerGroup
	}

	return d.orderKafkaConsumerGroup
}
