package payment

import (
	"context"
	repoModel "github.com/dfg007star/go_rocket/payment/internal/repository/model"
	"github.com/google/uuid"
)

func (r *repository) PayOrder(ctx context.Context, payment repoModel.Payment) (string, error) {

	r.mu.Lock()
	defer r.mu.Unlock()

	// network request example
	newUUID := uuid.NewString()
	payment.TransactionUuid = newUUID
	r.data = append(r.data, payment)

	return payment.TransactionUuid, nil
}
