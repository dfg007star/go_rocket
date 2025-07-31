package app

import (
	"context"
	"fmt"
	orderAPI "github.com/dfg007star/go_rocket/order/internal/api/order/v1"
	grpcClient "github.com/dfg007star/go_rocket/order/internal/client/grpc"
	inventoryServiceClient "github.com/dfg007star/go_rocket/order/internal/client/grpc/inventory/v1"
	paymentServiceClient "github.com/dfg007star/go_rocket/order/internal/client/grpc/payment/v1"
	"github.com/dfg007star/go_rocket/order/internal/config"
	"github.com/dfg007star/go_rocket/order/internal/repository"
	orderRepository "github.com/dfg007star/go_rocket/order/internal/repository/order"
	"github.com/dfg007star/go_rocket/order/internal/service"
	orderService "github.com/dfg007star/go_rocket/order/internal/service/order"
	orderV1 "github.com/dfg007star/go_rocket/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/payment/v1"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type diContainer struct {
	orderV1API      *orderV1.Server
	orderService    service.OrderService
	orderRepository repository.OrderRepository

	postgresClient  *pgx.Conn
	paymentClient   grpcClient.PaymentClient
	inventoryClient grpcClient.InventoryClient
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
		d.orderService = orderService.NewOrderService(d.OrderRepository(ctx), d.InventoryClient(ctx), d.PaymentClient(ctx))
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

		defer func() {
			cerr := con.Close(ctx)
			if cerr != nil {
				panic(fmt.Errorf("failed to close connection: %w", cerr))
			}
		}()

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

		client := inventoryV1.NewInventoryServiceClient(conn)
		d.inventoryClient = inventoryServiceClient.NewClient(client)
	}

	return d.inventoryClient
}
