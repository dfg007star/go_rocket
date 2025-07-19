package payment

import (
	"context"
	"fmt"

	"github.com/dfg007star/go_rocket/payment/internal/model"
	"github.com/dfg007star/go_rocket/payment/internal/repository/converter"
)

func (s *service) PayOrder(ctx context.Context, payment model.Payment) (string, error) {
	if err := s.validatePayment(payment); err != nil {
		return "", fmt.Errorf("%w: %w", err, model.ErrPayOrderModelValidation)
	}

	transactionUUID, err := s.paymentRepository.PayOrder(ctx, converter.PaymentToRepoModel(payment))
	if err != nil {
		return "", err
	}

	return transactionUUID, nil
}
