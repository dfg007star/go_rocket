package part

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	def "github.com/dfg007star/go_rocket/inventory/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ def.PartRepository = (*repository)(nil)

type repository struct {
	mu   sync.RWMutex
	data *mongo.Collection
}

func NewRepository(clientMongo *mongo.Client) *repository {
	dbName := os.Getenv("MONGO_INITDB_DATABASE")
	if dbName == "" {
		log.Fatal("MONGO_INITDB_DATABASE environment variable is not set")
	}

	db := clientMongo.Database(dbName)
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := data.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		panic(err)
	}

	return &repository{
		data: data,
	}
}
