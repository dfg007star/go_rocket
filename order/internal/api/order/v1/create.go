package v1

import (
	"context"
	"github.com/dfg007star/go_rocket/order/internal/model"
	orderV1 "github.com/dfg007star/go_rocket/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	order, err := a.orderService.Create(ctx, &model.OrderCreate{UserUuid: req.UserUUID, PartUuids: req.PartUuids})
	if err != nil {
		return nil, err
	}

	return &orderV1.CreateOrderResponse{
		OrderUUID:  order.OrderUuid,
		TotalPrice: order.TotalPrice,
	}, nil
}
