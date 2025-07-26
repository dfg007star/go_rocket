package repository

import (
	"context"

	"github.com/dfg007star/go_rocket/order/internal/model"
)

type OrderRepository interface {
	Create(ctx context.Context, userUuid string, parts []*model.Part) (*model.Order, error)
	Update(ctx context.Context, order *model.OrderUpdate) (*model.Order, error)
	Get(ctx context.Context, orderUuid string) (*model.Order, error)
}
