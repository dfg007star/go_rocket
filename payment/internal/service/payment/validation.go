package payment

import (
	"errors"
	"github.com/dfg007star/go_rocket/payment/internal/model"
)

func (s *service) validatePayment(payment model.Payment) error {
	if payment.OrderUuid == "" {
		return errors.New("order UUID is required")
	}

	if payment.UserUuid == "" {
		return errors.New("user UUID is required")
	}

	if payment.PaymentMethod < model.PAYMENT_METHOD_CARD ||
		payment.PaymentMethod > model.PAYMENT_METHOD_INVESTOR_MONEY {
		return errors.New("invalid or unspecified payment method")
	}

	return nil
}
