package converter

import (
	"github.com/dfg007star/go_rocket/order/internal/model"
	repoModel "github.com/dfg007star/go_rocket/order/internal/repository/model"
)

func RepoModelToOrder(order repoModel.Order) model.Order {
	return model.Order{
		OrderUuid:       order.OrderUuid,
		UserUuid:        order.UserUuid,
		PartUuids:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUuid: order.TransactionUuid,
		PaymentMethod:   (*model.PaymentMethod)(order.PaymentMethod),
		Status:          (model.Status)(order.Status),
	}
}
