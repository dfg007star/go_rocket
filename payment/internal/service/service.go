package service

import (
	"context"
	"github.com/dfg007star/go_rocket/payment/internal/model"
)

type PaymentService interface {
	PayOrder(ctx context.Context, payment model.Payment) (string, error)
}
