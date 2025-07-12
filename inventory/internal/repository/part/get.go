package part

import (
	"context"
	"github.com/dfg007star/go_rocket/inventory/internal/model"
	"github.com/dfg007star/go_rocket/inventory/internal/repository/converter"
)

func (r *repository) Get(ctx context.Context, uuid string) (model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, part := range r.data {
		if part.Uuid == uuid {
			return converter.RepoModelToPartModel(&part), nil
		}
	}

	return model.Part{}, model.ErrPartNotFound
}
