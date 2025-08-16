package service

import (
	"context"

	"github.com/dfg007star/go_rocket/notification/internal/model"
)

type OrderAssembledConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type OrderPaidConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type TelegramService interface {
	SendPaidNotification(ctx context.Context, event model.OrderPaidEvent) error
	SendAssembledNotification(ctx context.Context, event model.ShipAssembledEvent) error
}
