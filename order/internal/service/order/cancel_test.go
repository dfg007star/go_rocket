package order

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/dfg007star/go_rocket/order/internal/model"
)

func (s *ServiceSuite) TestCancelOrderSuccess() {
	orderUuid := gofakeit.UUID()

	paymentMethod := model.INVESTOR_MONEY
	order := &model.Order{
		OrderUuid:     orderUuid,
		UserUuid:      gofakeit.UUID(),
		PartUuids:     []string{gofakeit.UUID()},
		TotalPrice:    float32(gofakeit.Price(100, 1000)),
		PaymentMethod: &paymentMethod,
		Status:        model.PENDING_PAYMENT,
		CreatedAt:     gofakeit.Date(),
	}

	status := model.CANCELLED
	orderUpdateInfo := &model.OrderUpdate{
		OrderUuid: orderUuid,
		Status:    &status,
	}

	s.orderRepository.On("Get", s.ctx, orderUuid).Return(order, nil).Once()
	s.orderRepository.On("Update", s.ctx, orderUpdateInfo).Return(order, nil).Once()
	err := s.service.Cancel(s.ctx, orderUuid)
	s.NoError(err)
}

func (s *ServiceSuite) TestCancelOrderErr() {
	orderUuid := gofakeit.UUID()
	expectedErr := model.ErrOrderAlreadyPaid

	paymentMethod := model.INVESTOR_MONEY
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
	err := s.service.Cancel(s.ctx, orderUuid)
	s.Error(err)
	s.Equal(expectedErr, err)
}

func (s *ServiceSuite) TestCancelOrderConflictErr() {
	orderUuid := gofakeit.UUID()
	expectedErr := model.ErrOrderAlreadyCancelled

	paymentMethod := model.INVESTOR_MONEY
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
	err := s.service.Cancel(s.ctx, orderUuid)
	s.Error(err)
	s.Equal(expectedErr, err)
}

func (s *ServiceSuite) TestCancelOrderInternalErr() {
	orderUUID := gofakeit.UUID()
	expectedErr := model.ErrOrderInternalError

	s.orderRepository.On("Get", s.ctx, orderUUID).Return(&model.Order{}, expectedErr).Once()
	err := s.service.Cancel(s.ctx, orderUUID)
	s.Error(err)
	s.Equal(expectedErr, err)
}

func (s *ServiceSuite) TestCancelOrderNotFoundErr() {
	orderUUID := gofakeit.UUID()
	expectedErr := model.ErrOrderNotFound

	s.orderRepository.On("Get", s.ctx, orderUUID).Return(&model.Order{}, expectedErr).Once()
	err := s.service.Cancel(s.ctx, orderUUID)
	s.Error(err)
	s.Equal(expectedErr, err)
}

func (s *ServiceSuite) TestCancelOrderUpdateErr() {
	orderUuid := gofakeit.UUID()
	expectedErr := model.ErrOrderNotFound

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

	status := model.CANCELLED
	orderUpdate := &model.OrderUpdate{
		OrderUuid: orderUuid,
		Status:    &status,
	}

	s.orderRepository.On("Get", s.ctx, orderUuid).Return(order, nil).Once()
	s.orderRepository.On("Update", s.ctx, orderUpdate).Return(&model.Order{}, expectedErr).Once()
	err := s.service.Cancel(s.ctx, orderUuid)
	s.Error(err)
	s.Equal(expectedErr, err)
}
