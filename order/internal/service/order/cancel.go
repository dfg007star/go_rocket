package order

import (
	"context"

	"github.com/dfg007star/go_rocket/order/internal/model"
)

func (s *service) Cancel(ctx context.Context, orderUuid string) error {
	order, err := s.orderRepository.Get(ctx, orderUuid)
	if err != nil {
		return err
	}

	if order.Status == model.PAID {
		return model.ErrOrderAlreadyPaid
	}

	if order.Status == model.CANCELLED {
		return model.ErrOrderAlreadyCancelled
	}

	cancelledStatus := model.CANCELLED
	_, err = s.orderRepository.Update(ctx, model.OrderUpdate{
		OrderUuid: orderUuid,
		Status:    &cancelledStatus,
	})

	return err
}
