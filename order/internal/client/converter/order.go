package converter

import (
	"github.com/dfg007star/go_rocket/order/internal/model"
	generatedPaymentV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/payment/v1"
)

func PaymentMethodToProto(paymentMethod *model.PaymentMethod) generatedPaymentV1.PaymentMethod {
	switch *paymentMethod {
	case model.UNSPECIFIED:
		return generatedPaymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	case model.CARD:
		return generatedPaymentV1.PaymentMethod_PAYMENT_METHOD_CARD
	case model.SBP:
		return generatedPaymentV1.PaymentMethod_PAYMENT_METHOD_SBP
	case model.CREDIT_CARD:
		return generatedPaymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case model.INVESTOR_MONEY:
		return generatedPaymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return generatedPaymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}
