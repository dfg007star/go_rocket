package v1

import (
	"context"
	"github.com/dfg007star/go_rocket/payment/internal/converter"
	payment_v1 "github.com/dfg007star/go_rocket/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *payment_v1.PayOrderRequest) (*payment_v1.PayOrderResponse, error) {
	uuid, err := a.paymentService.PayOrder(ctx, converter.PayOrderRequestToPaymentModel(req))

	if err != nil {
		return nil, err
	}

	return &payment_v1.PayOrderResponse{
		TransactionUuid: uuid,
	}, nil
}
