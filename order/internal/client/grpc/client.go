package grpc

import (
	"context"

	"github.com/dfg007star/go_rocket/order/internal/model"
)

type InventoryClient interface {
	ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error)
	GetPart(ctx context.Context, uuid string) (model.Part, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, paymentMethod model.PaymentMethod, userUuid, orderUuid string) (string, error)
}
