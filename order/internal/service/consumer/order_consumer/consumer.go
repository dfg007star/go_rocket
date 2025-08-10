package order_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/dfg007star/go_rocket/order/internal/converter/kafka"
	def "github.com/dfg007star/go_rocket/order/internal/service"
	"github.com/dfg007star/go_rocket/platform/pkg/kafka"
	"github.com/dfg007star/go_rocket/platform/pkg/logger"
)

var _ def.ConsumerService = (*service)(nil)

type service struct {
	orderPaidConsumer     kafka.Consumer
	orderAssembledDecoder kafkaConverter.OrderAssembledDecoder
}

func NewService(orderPaidConsumer kafka.Consumer, orderAssembledDecoder kafkaConverter.OrderAssembledDecoder) *service {
	return &service{
		orderPaidConsumer:     orderPaidConsumer,
		orderAssembledDecoder: orderAssembledDecoder,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting Assembly service")

	err := s.orderPaidConsumer.Consume(ctx, s.OrderAssembledHandler)
	if err != nil {
		logger.Error(ctx, "Consume from ufo.recorded topic error", zap.Error(err))
		return err
	}

	return nil
}
