package v1

import (
	"context"
	orderV1 "github.com/dfg007star/go_rocket/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrderByUuid(ctx context.Context, params orderV1.CancelOrderByUuidParams) (orderV1.CancelOrderByUuidRes, error) {
	err := a.orderService.Cancel(ctx, params.OrderUUID)
	if err != nil {
		return nil, err
	}

	return &orderV1.CancelOrderByUuidNoContent{}, nil
}
