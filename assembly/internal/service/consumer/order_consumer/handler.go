package order_consumer

import (
	"context"
	"time"

	"github.com/dfg007star/go_rocket/assembly/internal/model"
	"github.com/dfg007star/go_rocket/platform/pkg/kafka/consumer"
	"github.com/dfg007star/go_rocket/platform/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
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
		logger.Info(ctx, "starting assembly", zap.String("order_uuid", event.OrderUuid))

		time.Sleep(10 * time.Second)

		assembledEvent := model.ShipAssembledEvent{
			EventUuid:    uuid.New().String(),
			OrderUuid:    event.OrderUuid,
			UserUuid:     event.UserUuid,
			BuildTimeSec: 10,
		}

		err := s.orderProducer.ShipAssembled(ctx, assembledEvent)
		if err != nil {
			logger.Error(ctx, "failed to publish ShipAssembled", zap.Error(err))
		}

		logger.Info(ctx, "assembly finished", zap.String("order_uuid", event.OrderUuid))
	}()

	return nil
}
