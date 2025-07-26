package v1

import (
	"context"

	converter "github.com/dfg007star/go_rocket/order/internal/converter"
	orderV1 "github.com/dfg007star/go_rocket/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	pm := orderV1.OrderDtoPaymentMethod(req.PaymentMethod)
	convertedPm := converter.ConvertPaymentMethodToModel(&pm)
	order, err := a.orderService.Pay(ctx, params.OrderUUID, &convertedPm)
	if err != nil {
		return nil, err
	}

	return &orderV1.PayOrderResponse{
		TransactionUUID: *order.TransactionUuid,
	}, nil
}
