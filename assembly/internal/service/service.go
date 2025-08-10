package service

import (
	"context"

	"github.com/dfg007star/go_rocket/assembly/internal/model"
)

type OrderProducerService interface {
	ShipAssembled(ctx context.Context, event model.ShipAssembledEvent) error
}

type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}
