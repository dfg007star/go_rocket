package order

import (
	"context"
	"github.com/dfg007star/go_rocket/order/internal/model"
)

func (s *service) Pay(ctx context.Context, orderUuid string, method model.PaymentMethod) (model.Order, error) {
	order, err := s.orderRepository.Get(ctx, orderUuid)
	if err != nil {
		return model.Order{}, err
	}

	if order.Status == model.PAID {
		return model.Order{}, model.ErrOrderAlreadyPaid
	}

	paymentMethod := method
	if order.PaymentMethod != nil {
		paymentMethod = *order.PaymentMethod
	}

	transactionUuid, err := s.paymentClient.PayOrder(ctx, paymentMethod, orderUuid, order.UserUuid)
	if err != nil {
		return model.Order{}, err
	}

	paidStatus := model.PAID
	updatedOrder, err := s.orderRepository.Update(ctx, model.OrderUpdate{
		OrderUuid:       orderUuid,
		TransactionUuid: &transactionUuid,
		Status:          &paidStatus,
	})
	if err != nil {
		return model.Order{}, err
	}

	return updatedOrder, nil
}
