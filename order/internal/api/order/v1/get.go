package v1

import (
	"context"

	"github.com/dfg007star/go_rocket/order/internal/converter"
	orderV1 "github.com/dfg007star/go_rocket/shared/pkg/openapi/order/v1"
)

func (a *api) OrderByUuid(ctx context.Context, params orderV1.OrderByUuidParams) (orderV1.OrderByUuidRes, error) {
	order, err := a.orderService.Get(ctx, params.OrderUUID)
	if err != nil {
		return nil, err
	}
	return converter.OrderModelToOrderDto(order), nil
}
