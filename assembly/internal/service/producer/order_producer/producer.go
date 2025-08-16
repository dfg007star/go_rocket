package order_producer

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/dfg007star/go_rocket/assembly/internal/model"
	"github.com/dfg007star/go_rocket/platform/pkg/kafka"
	"github.com/dfg007star/go_rocket/platform/pkg/logger"
	eventsV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/events/v1"
)

type service struct {
	orderProducer kafka.Producer
}

func NewService(orderProducer kafka.Producer) *service {
	return &service{
		orderProducer: orderProducer,
	}
}

func (p *service) ShipAssembled(ctx context.Context, event model.ShipAssembledEvent) error {
	msg := &eventsV1.ShipAssembled{
		EventUuid:    event.EventUuid,
		OrderUuid:    event.OrderUuid,
		UserUuid:     event.UserUuid,
		BuildTimeSec: event.BuildTimeSec,
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "failed to marshal UFORecorded", zap.Error(err))
		return err
	}

	err = p.orderProducer.Send(ctx, []byte(event.EventUuid), payload)
	if err != nil {
		logger.Error(ctx, "failed to publish ShipAssembled", zap.Error(err))
		return err
	}

	return nil
}
