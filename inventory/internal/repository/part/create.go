package part

import (
	"context"
	"github.com/dfg007star/go_rocket/inventory/internal/model"
	"github.com/dfg007star/go_rocket/inventory/internal/repository/converter"
)

func (r *repository) Create(ctx context.Context, part model.Part) (model.Part, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data = append(r.data, converter.PartModelToRepoModel(&part))

	return part, nil
}
