package part

import (
	"context"
	"fmt"
	"time"

	"github.com/dfg007star/go_rocket/inventory/internal/model"
	"github.com/dfg007star/go_rocket/inventory/internal/repository/converter"
	"github.com/google/uuid"
)

func (r *repository) Create(ctx context.Context, part *model.Part) (*model.Part, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	if part.Uuid == "" {
		part.CreatedAt = now
		part.Uuid = uuid.New().String()
	} else {
		return &model.Part{}, fmt.Errorf("part already exists: %s", part.Uuid)
	}
	part.UpdatedAt = now

	partMongo := converter.PartModelToRepoModel(part)

	_, err := r.data.InsertOne(ctx, partMongo)
	if err != nil {
		return nil, fmt.Errorf("failed to create or update part: %w", err)
	}

	result := converter.RepoModelToPartModel(partMongo)

	return result, nil
}
