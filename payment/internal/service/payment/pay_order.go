package payment

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/dfg007star/go_rocket/payment/internal/model"
	"github.com/dfg007star/go_rocket/payment/internal/repository/converter"
	"github.com/dfg007star/go_rocket/platform/pkg/tracing"
)

func (s *service) PayOrder(ctx context.Context, payment model.Payment) (string, error) {
	if err := s.validatePayment(payment); err != nil {
		return "", fmt.Errorf("%w: %w", err, model.ErrPayOrderModelValidation)
	}

	ctx, paymentSpan := tracing.StartSpan(ctx, "payment.pay_order",
		trace.WithAttributes(
			attribute.String("order.uuid", payment.OrderUuid),
			attribute.String("user.uuid", payment.UserUuid),
			attribute.String("payment_method", payment.PaymentMethod.String()),
		),
	)
	defer paymentSpan.End()

	transactionUUID, err := s.paymentRepository.PayOrder(ctx, converter.PaymentToRepoModel(payment))
	if err != nil {
		paymentSpan.RecordError(err)
		paymentSpan.End()
		return "", err
	}

	return transactionUUID, nil
}
