package repository

import (
	"context"

	"github.com/dfg007star/go_rocket/inventory/internal/model"
)

type InventoryRepository interface {
	Get(ctx context.Context, uuid string) (*model.Part, error)
	List(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error)
	Create(ctx context.Context, part *model.Part) (*model.Part, error)
	Update(ctx context.Context, part *model.Part) (*model.Part, error)
}
