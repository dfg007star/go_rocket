package app

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"

	inventoryV1API "github.com/dfg007star/go_rocket/inventory/internal/api/inventory/v1"
	"github.com/dfg007star/go_rocket/inventory/internal/config"
	"github.com/dfg007star/go_rocket/inventory/internal/model"
	"github.com/dfg007star/go_rocket/inventory/internal/repository"
	inventoryRepository "github.com/dfg007star/go_rocket/inventory/internal/repository/part"
	"github.com/dfg007star/go_rocket/inventory/internal/service"
	inventoryService "github.com/dfg007star/go_rocket/inventory/internal/service/part"
	"github.com/dfg007star/go_rocket/platform/pkg/closer"
	"github.com/dfg007star/go_rocket/platform/pkg/logger"
	inventoryV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/inventory/v1"
)

type diContainer struct {
	inventoryV1API inventoryV1.InventoryServiceServer

	inventoryService service.InventoryService

	inventoryRepository repository.InventoryRepository

	mongoDBClient *mongo.Client
	mongoDBHandle *mongo.Database
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) InventoryV1API(ctx context.Context) inventoryV1.InventoryServiceServer {
	if d.inventoryV1API == nil {
		d.inventoryV1API = inventoryV1API.NewAPI(d.InventoryService(ctx))
	}

	return d.inventoryV1API
}

func (d *diContainer) InventoryService(ctx context.Context) service.InventoryService {
	if d.inventoryService == nil {
		d.inventoryService = inventoryService.NewService(d.InventoryRepository(ctx))
	}

	return d.inventoryService
}

func (d *diContainer) InventoryRepository(ctx context.Context) repository.InventoryRepository {
	if d.inventoryRepository == nil {
		d.inventoryRepository = inventoryRepository.NewRepository(ctx, d.MongoDBHandle(ctx))
		// create part for test purpose
		part := model.Part{
			Name:          "Turbo Engine",
			Description:   "High-performance aircraft engine",
			Price:         125000.99,
			StockQuantity: 5,
			Category:      model.ENGINE,
			Dimensions: model.Dimensions{
				Length: 2.5,
				Width:  1.8,
				Height: 1.2,
				Weight: 350.5,
			},
			Manufacturer: model.Manufacturer{
				Name:    "AeroTech",
				Country: "USA",
				Website: "www.aerotech.example",
			},
			Tags:      []string{"engine", "turbo", "aircraft"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		_, err := d.inventoryRepository.Create(ctx, &part)
		if err != nil {
			logger.Error(ctx, "create inventory part error", zap.Error(err))
		}
	}

	return d.inventoryRepository
}

func (d *diContainer) MongoDBClient(ctx context.Context) *mongo.Client {
	if d.mongoDBClient == nil {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
		if err != nil {
			panic(fmt.Sprintf("failed to connect to MongoDB: %s\n", err.Error()))
		}

		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			panic(fmt.Sprintf("failed to ping MongoDB: %v\n", err))
		}

		closer.AddNamed("MongoDB client", func(ctx context.Context) error {
			return client.Disconnect(ctx)
		})

		d.mongoDBClient = client
	}

	return d.mongoDBClient
}

func (d *diContainer) MongoDBHandle(ctx context.Context) *mongo.Database {
	if d.mongoDBHandle == nil {
		d.mongoDBHandle = d.MongoDBClient(ctx).Database(config.AppConfig().Mongo.DatabaseName())
	}

	return d.mongoDBHandle
}
