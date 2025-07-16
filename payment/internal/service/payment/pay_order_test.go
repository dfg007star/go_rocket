package payment

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/dfg007star/go_rocket/payment/internal/model"
	"github.com/dfg007star/go_rocket/payment/internal/repository/converter"
	"github.com/stretchr/testify/require"
)

func (s *ServiceSuite) TestPayOrder() {
	var (
		orderUuid       = gofakeit.UUID()
		userUuid        = gofakeit.UUID()
		transactionUuid = gofakeit.UUID()
		paymentMethod   = model.PAYMENT_METHOD_CARD

		payment = model.Payment{
			OrderUuid:     orderUuid,
			UserUuid:      userUuid,
			PaymentMethod: paymentMethod,
		}

		repoPayment = converter.PaymentToRepoModel(payment)
	)

	s.paymentRepository.On("PayOrder", s.ctx, repoPayment).Return(transactionUuid, nil)

	resultUuid, err := s.service.PayOrder(s.ctx, payment)
	s.Require().NoError(err)
	s.Require().Equal(transactionUuid, resultUuid)
	s.paymentRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestPayOrderError() {
	s.Run("InvalidPaymentMethod", func() {
		payment := model.Payment{
			OrderUuid:     gofakeit.UUID(),
			UserUuid:      gofakeit.UUID(),
			PaymentMethod: model.PAYMENT_METHOD_UNSPECIFIED,
		}

		resultUuid, err := s.service.PayOrder(s.ctx, payment)
		require.Error(s.T(), err)
		require.Empty(s.T(), resultUuid)
		require.Contains(s.T(), err.Error(), "invalid order model")
	})

	s.Run("EmptyOrderUUID", func() {
		payment := model.Payment{
			OrderUuid:     "",
			UserUuid:      gofakeit.UUID(),
			PaymentMethod: model.PAYMENT_METHOD_CREDIT_CARD,
		}

		resultUuid, err := s.service.PayOrder(s.ctx, payment)
		require.Error(s.T(), err)
		require.Empty(s.T(), resultUuid)
		require.Contains(s.T(), err.Error(), "order UUID is required")
	})

	s.Run("EmptyUserUUID", func() {
		payment := model.Payment{
			OrderUuid:     gofakeit.UUID(),
			UserUuid:      "",
			PaymentMethod: model.PAYMENT_METHOD_INVESTOR_MONEY,
		}

		resultUuid, err := s.service.PayOrder(s.ctx, payment)
		require.Error(s.T(), err)
		require.Empty(s.T(), resultUuid)
		require.Contains(s.T(), err.Error(), "user UUID is required")
	})
}
