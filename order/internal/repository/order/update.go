package order

import (
	"context"
	"github.com/dfg007star/go_rocket/order/internal/model"
	"github.com/dfg007star/go_rocket/order/internal/repository/converter"
	repoModel "github.com/dfg007star/go_rocket/order/internal/repository/model"
)

func (r *repository) Update(ctx context.Context, orderUpdate model.OrderUpdate) (model.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, order := range r.data {
		if order.OrderUuid == orderUpdate.OrderUuid {
			if orderUpdate.TransactionUuid != nil {
				order.TransactionUuid = orderUpdate.TransactionUuid
			}

			if orderUpdate.PaymentMethod != nil {
				paymentMethod := (*repoModel.PaymentMethod)(orderUpdate.PaymentMethod)
				order.PaymentMethod = paymentMethod
			}

			if orderUpdate.Status != nil {
				order.Status = repoModel.Status(*orderUpdate.Status)
			}

			r.data[i] = order

			return converter.RepoModelToOrder(order), nil
		}
	}

	return model.Order{}, model.ErrOrderNotFound
}
