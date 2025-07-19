package order

import (
	"context"

	"github.com/dfg007star/go_rocket/order/internal/model"
	"github.com/dfg007star/go_rocket/order/internal/repository/converter"
)

func (r *repository) Get(ctx context.Context, orderUuid string) (model.Order, error) {
	for _, order := range r.data {
		if order.OrderUuid == orderUuid {
			return converter.RepoModelToOrder(order), nil
		}
	}

	return model.Order{}, model.ErrOrderNotFound
}
