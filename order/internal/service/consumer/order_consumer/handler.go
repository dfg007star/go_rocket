package order_consumer

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/dfg007star/go_rocket/order/internal/model"
	"github.com/dfg007star/go_rocket/platform/pkg/kafka/consumer"
	"github.com/dfg007star/go_rocket/platform/pkg/logger"
)

func (s *service) OrderAssembledHandler(ctx context.Context, msg consumer.Message) error {
	event, err := s.orderAssembledDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode ShipAssembled", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Processing message",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("event_uuid", event.EventUuid),
		zap.String("order_uuid", event.OrderUuid),
		zap.String("user_uuid", event.UserUuid),
		zap.Int64("build_time_sec", event.BuildTimeSec),
	)

	order, err := s.orderRepository.Get(ctx, event.OrderUuid)
	if err != nil {
		logger.Error(ctx, "Failed to get order", zap.Error(err))
		return err
	}

	if order.Status != model.PAID {
		err := errors.New("order status is not PAID")
		logger.Error(ctx, "Order status is not PAID", zap.Error(err))
		return err
	}

	completedStatus := model.COMPLETED
	_, err = s.orderRepository.Update(ctx, &model.OrderUpdate{
		OrderUuid: event.OrderUuid,
		Status:    &completedStatus,
	})
	if err != nil {
		logger.Error(ctx, "Failed to update order", zap.Error(err))
		return err
	}

	return nil
}
