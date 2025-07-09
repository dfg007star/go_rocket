package converter

import (
	"github.com/dfg007star/go_rocket/payment/internal/model"
	repoModel "github.com/dfg007star/go_rocket/payment/internal/repository/model"
)

func PaymentToRepoModel(payment model.Payment) repoModel.Payment {
	return repoModel.Payment{
		OrderUuid:     payment.OrderUuid,
		UserUuid:      payment.UserUuid,
		PaymentMethod: paymentMethodToRepoPaymentMethod(payment.PaymentMethod),
	}
}

func paymentMethodToRepoPaymentMethod(paymentMethod model.PaymentMethod) repoModel.PaymentMethod {
	switch paymentMethod {
	case model.PAYMENT_METHOD_CARD:
		return repoModel.PAYMENT_METHOD_CARD
	case model.PAYMENT_METHOD_SBP:
		return repoModel.PAYMENT_METHOD_SBP
	case model.PAYMENT_METHOD_CREDIT_CARD:
		return repoModel.PAYMENT_METHOD_CREDIT_CARD
	case model.PAYMENT_METHOD_INVESTOR_MONEY:
		return repoModel.PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return repoModel.PAYMENT_METHOD_UNSPECIFIED
	}
}
