package order

import (
	"context"
	"github.com/dfg007star/go_rocket/order/internal/model"
)

func (s *service) Create(ctx context.Context, orderCreate *model.OrderCreate) (model.Order, error) {
	parts, err := s.inventoryClient.ListParts(ctx, model.PartsFilter{Uuids: orderCreate.PartUuids})
	if err != nil {
		return model.Order{}, err
	}
	if len(orderCreate.PartUuids) != len(parts) {
		return model.Order{}, model.ErrNotAllPartsMatched
	}

	order, cerr := s.orderRepository.Create(ctx, orderCreate.UserUuid, parts)
	if cerr != nil {
		return model.Order{}, cerr
	}

	return order, nil
}
