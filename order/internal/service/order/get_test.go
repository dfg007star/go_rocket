package order

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/dfg007star/go_rocket/order/internal/model"
)

func (s *ServiceSuite) TestGetOrderSuccess() {
	orderUuid := gofakeit.UUID()

	paymentMethod := model.PaymentMethod(model.INVESTOR_MONEY)
	order := model.Order{
		OrderUuid:     orderUuid,
		UserUuid:      gofakeit.UUID(),
		PartUuids:     []string{gofakeit.UUID()},
		TotalPrice:    float32(gofakeit.Price(100, 1000)),
		PaymentMethod: &paymentMethod,
		Status:        model.PENDING_PAYMENT,
		CreatedAt:     gofakeit.Date(),
	}

	s.orderRepository.On("Get", s.ctx, orderUuid).Return(order, nil).Once()
	resp, err := s.service.Get(s.ctx, orderUuid)

	s.NoError(err)
	s.Equal(order, resp)
}

func (s *ServiceSuite) TestGetOrderNotFoundErr() {
	orderUuid := gofakeit.UUID()
	expectedErr := model.ErrOrderNotFound

	s.orderRepository.On("Get", s.ctx, orderUuid).Return(model.Order{}, expectedErr).Once()
	resp, err := s.service.Get(s.ctx, orderUuid)

	s.Error(err)
	s.Equal(expectedErr, err)
	s.Empty(resp)
}

func (s *ServiceSuite) TestGetOrderInternalErr() {
	orderUuid := gofakeit.UUID()
	expectedErr := model.ErrOrderInternalError

	s.orderRepository.On("Get", s.ctx, orderUuid).Return(model.Order{}, expectedErr).Once()
	resp, err := s.service.Get(s.ctx, orderUuid)

	s.Error(err)
	s.Equal(expectedErr, err)
	s.Empty(resp)
}
