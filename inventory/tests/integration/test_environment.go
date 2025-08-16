//go:build integration

package integration

import (
	"context"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	repoModel "github.com/dfg007star/go_rocket/inventory/internal/repository/model"
	inventoryV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/inventory/v1"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (env *TestEnvironment) InsertTestPart(ctx context.Context) (string, error) {
	partUuid := gofakeit.UUID()
	now := time.Now()

	partDoc := bson.M{
		"uuid":           partUuid,
		"name":           gofakeit.BeerName() + " " + gofakeit.CarModel(),
		"description":    gofakeit.HackerPhrase(),
		"price":          gofakeit.Float64Range(1000, 100000),
		"stock_quantity": 10,
		"category":       gofakeit.RandomInt([]int{int(repoModel.ENGINE), int(repoModel.FUEL), int(repoModel.PORTHOLE), int(repoModel.WING)}),
		"dimensions": bson.M{
			"length": gofakeit.Float64Range(0.5, 5.0),
			"width":  gofakeit.Float64Range(0.5, 3.0),
			"height": gofakeit.Float64Range(0.2, 2.0),
			"weight": gofakeit.Float64Range(5, 500),
		},
		"manufacturer": bson.M{
			"name":    gofakeit.Company(),
			"country": gofakeit.Country(),
			"website": "www." + gofakeit.DomainName(),
		},
		"tags": []string{
			gofakeit.Word(),
			gofakeit.Word(),
			gofakeit.Word(),
		},
		"created_at": now,
		"updated_at": now,
	}

	// Используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "inventory-service" // fallback значение
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(inventoryCollectionName).InsertOne(ctx, partDoc)
	if err != nil {
		return "", err
	}

	return partUuid, nil
}

func (env *TestEnvironment) InsertTestPartWithData(ctx context.Context, part *inventoryV1.Part) (string, error) {
	partUuid := gofakeit.UUID()

	partDoc := bson.M{
		"uuid":           partUuid,
		"name":           part.GetName(),
		"description":    part.GetDescription(),
		"price":          part.GetPrice(),
		"stock_quantity": part.GetStockQuantity(),
		"category":       part.GetCategory(),
		"dimensions":     part.GetDimensions(),
		"manufacturer":   part.GetManufacturer(),
		"tags":           part.GetTags(),
		"created_at":     part.GetCreatedAt(),
		"updated_at":     part.GetUpdatedAt(),
	}

	// Используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "inventory-service" // fallback значение
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(inventoryCollectionName).InsertOne(ctx, partDoc)
	if err != nil {
		return "", err
	}

	return partUuid, nil
}

func (env *TestEnvironment) GetTestPartInfo() *inventoryV1.Part {
	return &inventoryV1.Part{
		Uuid:          gofakeit.UUID(),
		Name:          gofakeit.BeerName(),
		Description:   gofakeit.HackerPhrase(),
		Price:         gofakeit.Float64Range(1, 1000),
		StockQuantity: 10,
		Category:      inventoryV1.Category_CATEGORY_FUEL,
		CreatedAt:     timestamppb.New(time.Now().Add(-2 * time.Hour)),
		UpdatedAt:     timestamppb.New(time.Now().Add(-1 * time.Hour)),
	}
}

func (env *TestEnvironment) ClearPartsCollection(ctx context.Context) error {
	// Используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "inventory-service" // fallback значение
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(inventoryCollectionName).DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}

	return nil
}
