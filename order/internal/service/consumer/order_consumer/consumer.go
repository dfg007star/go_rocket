package order_consumer

import (
	"context"

	kafkaConverter "github.com/dfg007star/go_rocket/order/internal/converter/kafka"
	"github.com/dfg007star/go_rocket/order/internal/repository"
	def "github.com/dfg007star/go_rocket/order/internal/service"
	"github.com/dfg007star/go_rocket/platform/pkg/kafka"
	"github.com/dfg007star/go_rocket/platform/pkg/logger"
	"go.uber.org/zap"
)

var _ def.ConsumerService = (*service)(nil)

type service struct {
	orderAssembledConsumer kafka.Consumer
	orderAssembledDecoder  kafkaConverter.OrderAssembledDecoder
	orderRepository        repository.OrderRepository
}

func NewService(orderAssembledConsumer kafka.Consumer, orderAssembledDecoder kafkaConverter.OrderAssembledDecoder, orderRepository repository.OrderRepository) *service {
	return &service{
		orderAssembledConsumer: orderAssembledConsumer,
		orderAssembledDecoder:  orderAssembledDecoder,
		orderRepository:        orderRepository,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting Order service")

	err := s.orderAssembledConsumer.Consume(ctx, s.OrderAssembledHandler)
	if err != nil {
		logger.Error(ctx, "Consume from ship.assembled topic error", zap.Error(err))
		return err
	}

	return nil
}
