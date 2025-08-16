package order_producer

import (
	"context"

	"github.com/dfg007star/go_rocket/order/internal/model"
	"github.com/dfg007star/go_rocket/platform/pkg/kafka"
	"github.com/dfg007star/go_rocket/platform/pkg/logger"
	eventsV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/events/v1"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type service struct {
	orderProducer kafka.Producer
}

func NewService(orderProducer kafka.Producer) *service {
	return &service{
		orderProducer: orderProducer,
	}
}

func (p *service) OrderPaid(ctx context.Context, event model.OrderPaidEvent) error {
	msg := &eventsV1.OrderPaid{
		EventUuid:       event.EventUuid,
		OrderUuid:       event.OrderUuid,
		UserUuid:        event.UserUuid,
		PaymentMethod:   event.PaymentMethod,
		TransactionUuid: event.TransactionUuid,
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "failed to marshal UFORecorded", zap.Error(err))
		return err
	}

	err = p.orderProducer.Send(ctx, []byte(event.EventUuid), payload)
	if err != nil {
		logger.Error(ctx, "failed to publish OrderPaid", zap.Error(err))
		return err
	}

	return nil
}
