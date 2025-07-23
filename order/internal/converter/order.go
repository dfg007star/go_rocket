package converter

import (
	"github.com/dfg007star/go_rocket/order/internal/model"
	orderV1 "github.com/dfg007star/go_rocket/shared/pkg/openapi/order/v1"
)

func OrderModelToOrderDto(order *model.Order) *orderV1.OrderDto {
	var tu string
	if order.TransactionUuid != nil {
		tu = *order.TransactionUuid
	}
	return &orderV1.OrderDto{
		OrderUUID:       order.OrderUuid,
		UserUUID:        order.UserUuid,
		PartUuids:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: orderV1.OptString{Value: tu},
		PaymentMethod:   orderV1.OptOrderDtoPaymentMethod{Value: convertPaymentMethod(order.PaymentMethod)},
		Status:          convertStatus(order.Status),
		CreatedAt:       order.CreatedAt,
	}
}

func convertStatus(status model.Status) orderV1.OrderDtoStatus {
	switch status {
	case 0:
		return orderV1.OrderDtoStatusPENDINGPAYMENT
	case 1:
		return orderV1.OrderDtoStatusPAID
	case 2:
		return orderV1.OrderDtoStatusCANCELLED
	default:
		return ""
	}
}

func convertPaymentMethod(paymentMethod *model.PaymentMethod) orderV1.OrderDtoPaymentMethod {
	if paymentMethod == nil {
		return ""
	}

	switch *paymentMethod {
	case model.CARD:
		return orderV1.OrderDtoPaymentMethodPAYMENTMETHODCARD
	case model.SBP:
		return orderV1.OrderDtoPaymentMethodPAYMENTMETHODSBP
	case model.CREDIT_CARD:
		return orderV1.OrderDtoPaymentMethodPAYMENTMETHODCARD
	case model.INVESTOR_MONEY:
		return orderV1.OrderDtoPaymentMethodPAYMENTMETHODINVESTORMONEY
	default:
		return ""
	}
}

func ConvertPaymentMethodToModel(paymentMethod *orderV1.OrderDtoPaymentMethod) model.PaymentMethod {
	switch *paymentMethod {
	case orderV1.OrderDtoPaymentMethodPAYMENTMETHODCARD:
		return model.CARD
	case orderV1.OrderDtoPaymentMethodPAYMENTMETHODSBP:
		return model.SBP
	case orderV1.OrderDtoPaymentMethodPAYMENTMETHODCREDITCARD:
		return model.CREDIT_CARD
	case orderV1.OrderDtoPaymentMethodPAYMENTMETHODINVESTORMONEY:
		return model.INVESTOR_MONEY
	default:
		return model.UNSPECIFIED
	}
}
