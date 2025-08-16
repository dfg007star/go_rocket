package order_consumer

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/dfg007star/go_rocket/notification/internal/model"
	"github.com/dfg007star/go_rocket/platform/pkg/kafka/consumer"
	"github.com/dfg007star/go_rocket/platform/pkg/logger"
)

func (s *service) OrderPaidHandler(ctx context.Context, msg consumer.Message) error {
	event, err := s.orderPaidDecoder.PaidDecode(msg.Value)
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

	paidEvent := model.OrderPaidEvent{
		EventUuid:       uuid.New().String(),
		OrderUuid:       event.OrderUuid,
		UserUuid:        event.UserUuid,
		TransactionUuid: event.TransactionUuid,
		PaymentMethod:   event.PaymentMethod,
	}

	err = s.telegramService.SendPaidNotification(ctx, paidEvent)
	if err != nil {
		logger.Error(ctx, "Failed to send paid notification", zap.Error(err))
	}

	return nil
}
