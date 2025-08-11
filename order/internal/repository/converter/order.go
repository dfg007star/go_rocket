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

func PaymentMethodFromString(str string) repoModel.PaymentMethod {
	switch str {
	case "CARD":
		return repoModel.CARD
	case "SBP":
		return repoModel.SBP
	case "CREDIT_CARD":
		return repoModel.CREDIT_CARD
	case "INVESTOR_MONEY":
		return repoModel.INVESTOR_MONEY
	default:
		return repoModel.UNSPECIFIED
	}
}
