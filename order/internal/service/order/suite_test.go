package order

import (
	"context"
	client "github.com/dfg007star/go_rocket/order/internal/client/grpc/mocks"
	"github.com/dfg007star/go_rocket/order/internal/repository/mocks"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	orderRepository *mocks.OrderRepository
	inventoryClient *client.InventoryClient
	paymentClient   *client.PaymentClient

	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.orderRepository = mocks.NewOrderRepository(s.T())
	s.inventoryClient = client.NewInventoryClient(s.T())
	s.paymentClient = client.NewPaymentClient(s.T())

	s.service = NewOrderService(
		s.orderRepository,
		s.inventoryClient,
		s.paymentClient,
	)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
