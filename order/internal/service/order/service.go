package order

import (
	"github.com/dfg007star/go_rocket/order/internal/client/grpc"
	"github.com/dfg007star/go_rocket/order/internal/repository"
	def "github.com/dfg007star/go_rocket/order/internal/service"
)

var _ def.OrderService = (*service)(nil)

type service struct {
	orderRepository repository.OrderRepository

	inventoryClient grpc.InventoryClient
	paymentClient   grpc.PaymentClient

	orderProducerService def.OrderProducerService
}

func NewOrderService(
	orderRepository repository.OrderRepository,
	inventoryClient grpc.InventoryClient,
	paymentClient grpc.PaymentClient,
	orderProducerService def.OrderProducerService,
) *service {
	return &service{
		orderRepository:      orderRepository,
		inventoryClient:      inventoryClient,
		paymentClient:        paymentClient,
		orderProducerService: orderProducerService,
	}
}
