package order_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/dfg007star/go_rocket/notification/internal/converter/kafka"
	def "github.com/dfg007star/go_rocket/notification/internal/service"
	"github.com/dfg007star/go_rocket/platform/pkg/kafka"
	"github.com/dfg007star/go_rocket/platform/pkg/logger"
)

var _ def.OrderAssembledConsumerService = (*service)(nil)

type service struct {
	orderAssembledConsumer kafka.Consumer
	orderAssembledDecoder  kafkaConverter.OrderAssembledDecoder
}

func NewService(orderAssembledConsumer kafka.Consumer, orderAssembledDecoder kafkaConverter.OrderAssembledDecoder) *service {
	return &service{
		orderAssembledConsumer: orderAssembledConsumer,
		orderAssembledDecoder:  orderAssembledDecoder,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting Assembly service")

	err := s.orderAssembledConsumer.Consume(ctx, s.OrderAssembledHandler)
	if err != nil {
		logger.Error(ctx, "Consume from ufo.recorded topic error", zap.Error(err))
		return err
	}

	return nil
}
