package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderAPI "github.com/dfg007star/go_rocket/order/internal/api/order/v1"
	inventoryServiceClient "github.com/dfg007star/go_rocket/order/internal/client/grpc/inventory/v1"
	paymentServiceClient "github.com/dfg007star/go_rocket/order/internal/client/grpc/payment/v1"
	"github.com/dfg007star/go_rocket/order/internal/config"
	orderRepository "github.com/dfg007star/go_rocket/order/internal/repository/order"
	orderService "github.com/dfg007star/go_rocket/order/internal/service/order"
	orderV1 "github.com/dfg007star/go_rocket/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/payment/v1"
)

const configPath = "../deploy/compose/order/.env"

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}
	ctx := context.Background()

	inventoryClient, _, err := newInventoryClient()
	if err != nil {
		panic(err)
	}
	paymentClient, _, err := newPaymentClient()
	if err != nil {
		panic(err)
	}

	con, err := pgx.Connect(ctx, config.AppConfig().Postgres.URI())
	if err != nil {
		log.Printf("failed to connect to database: %v\n", err)
		return
	}
	defer func() {
		cerr := con.Close(ctx)
		if cerr != nil {
			log.Printf("failed to close connection: %v\n", cerr)
		}
	}()

	err = con.Ping(ctx)
	if err != nil {
		log.Printf("База данных недоступна: %v\n", err)
		return
	}

	// init app
	repo := orderRepository.NewRepository(con)
	service := orderService.NewOrderService(
		repo,
		inventoryServiceClient.NewClient(inventoryClient),
		paymentServiceClient.NewClient(paymentClient),
	)
	api := orderAPI.NewApi(service)

	orderServer, err := orderV1.NewServer(api)
	if err != nil {
		panic(err)
	}

	// Инициализируем роутер Chi
	r := chi.NewRouter()

	// Добавляем middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	// Монтируем обработчики OpenAPI
	r.Mount("/", orderServer)

	server := &http.Server{
		Addr:              config.AppConfig().OrderHTTP.Address(),
		Handler:           r,
		ReadHeaderTimeout: config.AppConfig().OrderHTTP.ReadTimeout(),
	}

	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", config.AppConfig().OrderHTTP.Address())
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}

func newInventoryClient() (inventoryV1.InventoryServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		config.AppConfig().InventoryGRPC.Address(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, err
	}

	client := inventoryV1.NewInventoryServiceClient(conn)
	return client, conn, nil
}

func newPaymentClient() (paymentV1.PaymentServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		config.AppConfig().PaymentGRPC.Address(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, err
	}

	client := paymentV1.NewPaymentServiceClient(conn)
	return client, conn, nil
}
