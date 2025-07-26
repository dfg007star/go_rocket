package service

import (
	"context"

	"github.com/dfg007star/go_rocket/order/internal/model"
)

type OrderService interface {
	Get(ctx context.Context, orderUuid string) (*model.Order, error)
	Create(ctx context.Context, order *model.OrderCreate) (*model.Order, error)
	Cancel(ctx context.Context, orderUuid string) error
	Pay(ctx context.Context, orderUuid string, method *model.PaymentMethod) (*model.Order, error)
}
