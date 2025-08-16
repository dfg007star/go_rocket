package order

import (
	"errors"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/dfg007star/go_rocket/order/internal/model"
)

func (s *ServiceSuite) TestPayOrderSuccess() {
	orderUuid := gofakeit.UUID()
	paymentMethod := model.SBP
	transactionUuid := gofakeit.UUID()

	order := &model.Order{
		OrderUuid:     orderUuid,
		UserUuid:      gofakeit.UUID(),
		PartUuids:     []string{gofakeit.UUID()},
		TotalPrice:    float32(gofakeit.Price(100, 1000)),
		PaymentMethod: nil,
		Status:        model.PENDING_PAYMENT,
		CreatedAt:     gofakeit.Date(),
	}

	updatedStatus := model.PAID
	orderInfo := &model.OrderUpdate{
		OrderUuid:       orderUuid,
		TransactionUuid: &transactionUuid,
		Status:          &updatedStatus,
	}

	updatedOrder := *order
	updatedOrder.Status = updatedStatus
	updatedOrder.TransactionUuid = &transactionUuid

	s.orderRepository.On("Get", s.ctx, orderUuid).Return(order, nil).Once()
	s.paymentClient.On("PayOrder", s.ctx, &paymentMethod, order.OrderUuid, order.UserUuid).Return(transactionUuid, nil).Once()
	s.orderRepository.On("Update", s.ctx, orderInfo).Return(&updatedOrder, nil).Once()

	resp, err := s.service.Pay(s.ctx, orderUuid, &paymentMethod)
	s.NoError(err)
	s.Equal(transactionUuid, *resp.TransactionUuid)
}

func (s *ServiceSuite) TestPayOrderErrGetOrder() {
	orderUuid := gofakeit.UUID()
	paymentMethod := model.SBP
	expectedErr := model.ErrOrderNotFound

	s.orderRepository.On("Get", s.ctx, orderUuid).Return(&model.Order{}, expectedErr).Once()

	resp, err := s.service.Pay(s.ctx, orderUuid, &paymentMethod)
	s.Error(err)
	s.Equal(err, expectedErr)
	s.Empty(resp)
}

func (s *ServiceSuite) TestPayOrderErr() {
	orderUuid := gofakeit.UUID()
	paymentMethod := model.SBP

	order := &model.Order{
		OrderUuid:     orderUuid,
		UserUuid:      gofakeit.UUID(),
		PartUuids:     []string{gofakeit.UUID()},
		TotalPrice:    float32(gofakeit.Price(100, 1000)),
		PaymentMethod: &paymentMethod,
		Status:        model.PENDING_PAYMENT,
		CreatedAt:     gofakeit.Date(),
	}

	s.orderRepository.On("Get", s.ctx, orderUuid).Return(order, nil).Once()
	s.paymentClient.On("PayOrder", s.ctx, &paymentMethod, order.OrderUuid, order.UserUuid).Return("", gofakeit.Error()).Once()

	resp, err := s.service.Pay(s.ctx, orderUuid, &paymentMethod)
	s.Error(err)
	s.Empty(resp)
}

func (s *ServiceSuite) TestPayOrderInternalErr() {
	orderUuid := gofakeit.UUID()
	paymentMethod := model.CARD
	expectedErr := errors.New("order UUID is required")

	order := &model.Order{
		OrderUuid:     "",
		UserUuid:      gofakeit.UUID(),
		PartUuids:     []string{gofakeit.UUID()},
		TotalPrice:    float32(gofakeit.Price(100, 1000)),
		PaymentMethod: &paymentMethod,
		Status:        model.PENDING_PAYMENT,
		CreatedAt:     gofakeit.Date(),
	}

	s.orderRepository.On("Get", s.ctx, orderUuid).Return(order, nil).Once()
	s.paymentClient.On("PayOrder", s.ctx, &paymentMethod, order.OrderUuid, order.UserUuid).Return("", expectedErr).Once()

	resp, err := s.service.Pay(s.ctx, orderUuid, &paymentMethod)
	s.Error(err)
	s.Equal(err, expectedErr)
	s.Empty(resp)
}

func (s *ServiceSuite) TestPayOrderConflictOrderStatusPaidErr() {
	orderUuid := gofakeit.UUID()
	paymentMethod := model.SBP
	expectedErr := model.ErrOrderAlreadyPaid

	order := &model.Order{
		OrderUuid:     orderUuid,
		UserUuid:      gofakeit.UUID(),
		PartUuids:     []string{gofakeit.UUID()},
		TotalPrice:    float32(gofakeit.Price(100, 1000)),
		PaymentMethod: &paymentMethod,
		Status:        model.PAID,
		CreatedAt:     gofakeit.Date(),
	}

	s.orderRepository.On("Get", s.ctx, orderUuid).Return(order, nil).Once()

	resp, err := s.service.Pay(s.ctx, orderUuid, &paymentMethod)
	s.Error(err)
	s.Equal(err, expectedErr)
	s.Empty(resp)
}

func (s *ServiceSuite) TestPayOrderConflictOrderStatusCanceledErr() {
	orderUuid := gofakeit.UUID()
	paymentMethod := model.SBP
	expectedErr := model.ErrOrderAlreadyCancelled

	order := &model.Order{
		OrderUuid:     orderUuid,
		UserUuid:      gofakeit.UUID(),
		PartUuids:     []string{gofakeit.UUID()},
		TotalPrice:    float32(gofakeit.Price(100, 1000)),
		PaymentMethod: &paymentMethod,
		Status:        model.CANCELLED,
		CreatedAt:     gofakeit.Date(),
	}

	s.orderRepository.On("Get", s.ctx, orderUuid).Return(order, nil).Once()

	resp, err := s.service.Pay(s.ctx, orderUuid, &paymentMethod)
	s.Error(err)
	s.Equal(err, expectedErr)
	s.Empty(resp)
}
