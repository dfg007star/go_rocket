package service

import (
	"context"
)

type OrderAssembledConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type OrderPaidConsumerService interface {
	RunConsumer(ctx context.Context) error
}
