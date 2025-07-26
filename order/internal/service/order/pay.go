package order

import (
	"context"

	"github.com/dfg007star/go_rocket/order/internal/model"
)

func (s *service) Pay(ctx context.Context, orderUuid string, method *model.PaymentMethod) (*model.Order, error) {
	order, err := s.orderRepository.Get(ctx, orderUuid)
	if err != nil {
		return nil, err
	}

	if order.Status == model.PAID {
		return nil, model.ErrOrderAlreadyPaid
	}

	if order.Status == model.CANCELLED {
		return nil, model.ErrOrderAlreadyCancelled
	}

	paymentMethod := method
	if order.PaymentMethod != nil {
		paymentMethod = order.PaymentMethod
	}

	transactionUuid, err := s.paymentClient.PayOrder(ctx, paymentMethod, order.OrderUuid, order.UserUuid)
	if err != nil {
		return nil, err
	}

	paidStatus := model.PAID
	updatedOrder, err := s.orderRepository.Update(ctx, &model.OrderUpdate{
		OrderUuid:       order.OrderUuid,
		TransactionUuid: &transactionUuid,
		Status:          &paidStatus,
	})
	if err != nil {
		return nil, err
	}

	return updatedOrder, nil
}
