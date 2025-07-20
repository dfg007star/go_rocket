package part

import (
	"context"
	"github.com/dfg007star/go_rocket/inventory/internal/model"
	"github.com/dfg007star/go_rocket/inventory/internal/repository/converter"
	repoModel "github.com/dfg007star/go_rocket/inventory/internal/repository/model"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *repository) Get(ctx context.Context, uuid string) (*model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var part repoModel.Part

	err := r.data.FindOne(ctx, bson.M{"uuid": uuid}).Decode(&part)
	if err != nil {
		return nil, err
	}

	result := converter.RepoModelToPartModel(&part)

	return result, nil
}
