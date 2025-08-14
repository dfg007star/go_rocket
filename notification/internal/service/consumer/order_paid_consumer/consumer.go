package order_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/dfg007star/go_rocket/notification/internal/converter/kafka"
	def "github.com/dfg007star/go_rocket/notification/internal/service"
	"github.com/dfg007star/go_rocket/platform/pkg/kafka"
	"github.com/dfg007star/go_rocket/platform/pkg/logger"
)

var _ def.OrderPaidConsumerService = (*service)(nil)

type service struct {
	orderPaidConsumer kafka.Consumer
	orderPaidDecoder  kafkaConverter.OrderPaidDecoder
}

func NewService(orderPaidConsumer kafka.Consumer, orderPaidDecoder kafkaConverter.OrderPaidDecoder) *service {
	return &service{
		orderPaidConsumer: orderPaidConsumer,
		orderPaidDecoder:  orderPaidDecoder,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting Assembly service")

	err := s.orderPaidConsumer.Consume(ctx, s.OrderPaidHandler)
	if err != nil {
		logger.Error(ctx, "Consume from ufo.recorded topic error", zap.Error(err))
		return err
	}

	return nil
}
