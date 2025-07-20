package order

import (
	"context"

	"github.com/dfg007star/go_rocket/order/internal/model"
)

func (s *service) Get(ctx context.Context, orderUuid string) (model.Order, error) {
	order, err := s.orderRepository.Get(ctx, orderUuid)
	if err != nil {
		return model.Order{}, err
	}
	return order, nil
}
