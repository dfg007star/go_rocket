package converter

import (
	"github.com/dfg007star/go_rocket/payment/internal/model"
	payment_v1 "github.com/dfg007star/go_rocket/shared/pkg/proto/payment/v1"
)

func PayOrderRequestToPaymentModel(req *payment_v1.PayOrderRequest) model.Payment {
	return model.Payment{
		OrderUuid:     req.OrderUuid,
		UserUuid:      req.UserUuid,
		PaymentMethod: paymentMethodRequestToPaymentMethod(req.PaymentMethod),
	}
}

func paymentMethodRequestToPaymentMethod(paymentMethod payment_v1.PaymentMethod) model.PaymentMethod {
	switch paymentMethod {
	case payment_v1.PaymentMethod_PAYMENT_METHOD_CARD:
		return model.PAYMENT_METHOD_CARD
	case payment_v1.PaymentMethod_PAYMENT_METHOD_SBP:
		return model.PAYMENT_METHOD_SBP
	case payment_v1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return model.PAYMENT_METHOD_CREDIT_CARD
	case payment_v1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return model.PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return model.PAYMENT_METHOD_UNSPECIFIED
	}
}
