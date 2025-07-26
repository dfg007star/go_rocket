package part

import (
	"context"
	"fmt"
	"time"

	"github.com/dfg007star/go_rocket/inventory/internal/model"
	"github.com/dfg007star/go_rocket/inventory/internal/repository/converter"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *repository) Update(ctx context.Context, part *model.Part) (*model.Part, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	if part.Uuid == "" {
		part.Uuid = uuid.New().String()
		part.CreatedAt = now
	}
	part.UpdatedAt = now

	partMongo := converter.PartModelToRepoModel(part)

	filter := bson.M{"uuid": part.Uuid}

	update := bson.M{
		"$set": partMongo,
		"$setOnInsert": bson.M{
			"created_at": now,
		},
	}

	opts := options.Update().SetUpsert(true)

	_, err := r.data.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create or update part: %w", err)
	}

	result := converter.RepoModelToPartModel(partMongo)

	return result, nil
}
