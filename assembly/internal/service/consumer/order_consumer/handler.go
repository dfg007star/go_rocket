package order_consumer

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	assemblyMetrics "github.com/dfg007star/go_rocket/assembly/internal/metrics"
	"github.com/dfg007star/go_rocket/assembly/internal/model"
	"github.com/dfg007star/go_rocket/platform/pkg/kafka/consumer"
	"github.com/dfg007star/go_rocket/platform/pkg/logger"
)

func (s *service) OrderPaidHandler(ctx context.Context, msg consumer.Message) error {
	event, err := s.orderPaidDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderPaid", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Processing message",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("event_uuid", event.EventUuid),
		zap.String("order_uuid", event.OrderUuid),
		zap.String("user_uuid", event.UserUuid),
		zap.String("transaction_uuid", event.TransactionUuid),
		zap.String("payment_method", event.PaymentMethod),
	)

	go func() {
		start := time.Now()
		defer func() {
			duration := time.Since(start)
			assemblyMetrics.AssemblyDuration.Record(ctx, duration.Seconds())
		}()

		logger.Info(ctx, "Starting assembly", zap.String("order_uuid", event.OrderUuid))

		select {
		case <-time.After(10 * time.Second):
			assembledEvent := model.ShipAssembledEvent{
				EventUuid:    uuid.New().String(),
				OrderUuid:    event.OrderUuid,
				UserUuid:     event.UserUuid,
				BuildTimeSec: 10,
			}

			if err := s.orderProducer.ShipAssembled(ctx, assembledEvent); err != nil {
				logger.Error(ctx, "failed to publish ShipAssembled", zap.Error(err))
				return
			}

			logger.Info(ctx, "assembly finished", zap.String("order_uuid", event.OrderUuid))

		case <-ctx.Done():
			logger.Warn(ctx, "assembly canceled",
				zap.String("order_uuid", event.OrderUuid),
				zap.Error(ctx.Err()),
			)
			return
		}
	}()

	return nil
}
