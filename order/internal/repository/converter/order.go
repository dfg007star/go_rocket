package converter

import (
	"github.com/dfg007star/go_rocket/order/internal/model"
	repoModel "github.com/dfg007star/go_rocket/order/internal/repository/model"
)

func RepoModelToOrder(order *repoModel.Order) *model.Order {
	return &model.Order{
		OrderUuid:       order.OrderUuid,
		UserUuid:        order.UserUuid,
		PartUuids:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUuid: order.TransactionUuid,
		PaymentMethod:   (*model.PaymentMethod)(order.PaymentMethod),
		Status:          (model.Status)(order.Status),
	}
}

func OrderUpdateToRepoOrderUpdate(order *model.OrderUpdate) *repoModel.OrderUpdate {
	return &repoModel.OrderUpdate{
		OrderUuid:       order.OrderUuid,
		TransactionUuid: order.TransactionUuid,
		PaymentMethod:   (*repoModel.PaymentMethod)(order.PaymentMethod),
		Status:          (*repoModel.Status)(order.Status),
	}
}
