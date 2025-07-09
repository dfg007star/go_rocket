package payment

import (
	"context"
	"github.com/dfg007star/go_rocket/payment/internal/model"
	"github.com/dfg007star/go_rocket/payment/internal/repository/converter"
)

func (s *service) PayOrder(ctx context.Context, payment model.Payment) (string, error) {
	//create validation of model Payment
	//if err := val.ValidateStruct(request); err != nil {
	//	return "", fmt.Errorf("%w: %w", err, errors.ErrPayOrderModelValidationError)
	//}

	transactionUUID, err := s.paymentRepository.PayOrder(ctx, converter.PaymentToRepoModel(payment))
	if err != nil {
		return "", err
	}

	return transactionUUID, nil
}
