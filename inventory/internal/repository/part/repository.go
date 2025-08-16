package part

import (
	"context"

	def "github.com/dfg007star/go_rocket/inventory/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ def.InventoryRepository = (*repository)(nil)

type repository struct {
	data *mongo.Collection
}

func NewRepository(ctx context.Context, db *mongo.Database) *repository {
	data := db.Collection("parts")

	indexModels := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "uuid", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "name", Value: 1},
			},
		},
	}

	_, err := data.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		panic(err)
	}

	return &repository{
		data: data,
	}
}
