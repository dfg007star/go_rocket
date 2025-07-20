package repository

import (
	"context"

	repoModel "github.com/dfg007star/go_rocket/payment/internal/repository/model"
)

type PaymentRepository interface {
	PayOrder(ctx context.Context, payment repoModel.Payment) (string, error)
}
