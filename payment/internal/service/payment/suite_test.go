package payment

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/dfg007star/go_rocket/payment/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	paymentRepository *mocks.PaymentRepository

	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.paymentRepository = mocks.NewPaymentRepository(s.T())

	s.service = NewService(
		s.paymentRepository,
	)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
