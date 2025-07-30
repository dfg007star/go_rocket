package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	inventoryAPI "github.com/dfg007star/go_rocket/inventory/internal/api/inventory/v1"
	"github.com/dfg007star/go_rocket/inventory/internal/config"
	"github.com/dfg007star/go_rocket/inventory/internal/model"
	inventoryRepository "github.com/dfg007star/go_rocket/inventory/internal/repository/part"
	inventoryService "github.com/dfg007star/go_rocket/inventory/internal/service/part"
	inventoryV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/inventory/v1"
)

const configPath = "../deploy/compose/inventory/.env"

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("error loading config: %w", err))
	}

	ctx := context.Background()

	clientMongo, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
	if err != nil {
		log.Printf("failed to connect to database: %v\n", err)
		return
	}
	defer func() {
		cerr := clientMongo.Disconnect(ctx)
		if cerr != nil {
			log.Printf("failed to disconnect: %v\n", cerr)
		}
	}()

	err = clientMongo.Ping(ctx, nil)
	if err != nil {
		log.Printf("failed to ping database: %v\n", err)
		return
	}

	lis, err := net.Listen("tcp", config.AppConfig().InventoryGRPC.Address())
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}

	defer func() {
		if err := lis.Close(); err != nil {
			log.Printf("failed to close listener: %v\n", err)
		}
	}()

	repo := inventoryRepository.NewRepository(clientMongo)

	temperature_data := "-30°C to 80°C"
	sampleParts := []model.Part{
		{
			Name:          "TurboJet X-200 Engine",
			Description:   "High-performance turbofan engine for commercial aircraft",
			Price:         4250.99,
			StockQuantity: 15,
			Category:      model.ENGINE,
			Dimensions: model.Dimensions{
				Length: 3.2,
				Width:  1.8,
				Height: 2.1,
				Weight: 450.5,
			},
			Manufacturer: model.Manufacturer{
				Name:    "Boeing Aerospace",
				Country: "USA",
				Website: "www.boeing.com",
			},
			Tags: []string{"commercial", "high-efficiency"},
			Metadata: map[string]model.Value{
				"temperature_range": {StringValue: &temperature_data},
			},
			CreatedAt: time.Date(2023, 5, 10, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			Name:          "AeroFlap Wing Component",
			Description:   "Durable wing flap for mid-size aircraft",
			Price:         1250.50,
			StockQuantity: 32,
			Category:      model.WING,
			Dimensions: model.Dimensions{
				Length: 2.5,
				Width:  0.8,
				Height: 0.3,
				Weight: 85.2,
			},
			Manufacturer: model.Manufacturer{
				Name:    "Airbus Components",
				Country: "France",
				Website: "www.airbus.com",
			},
			Tags: []string{"lightweight", "composite"},
			Metadata: map[string]model.Value{
				"temperature_range": {StringValue: &temperature_data},
			},
			CreatedAt: time.Date(2024, 2, 20, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 2, 20, 0, 0, 0, 0, time.UTC),
		},
	}

	service := inventoryService.NewService(repo)
	api := inventoryAPI.NewApi(service)

	for d := range sampleParts {
		_, err := repo.Create(ctx, &sampleParts[d])
		if err != nil {
			continue
		}
	}

	// Создаем gRPC сервер
	s := grpc.NewServer()
	inventoryV1.RegisterInventoryServiceServer(s, api)

	// Включаем рефлексию для отладки
	reflection.Register(s)

	go func() {
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
