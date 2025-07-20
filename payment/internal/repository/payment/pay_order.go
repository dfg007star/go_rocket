package payment

import (
	"context"
	"log"

	"github.com/google/uuid"

	repoModel "github.com/dfg007star/go_rocket/payment/internal/repository/model"
)

func (r *repository) PayOrder(ctx context.Context, payment repoModel.Payment) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// network request example
	newUUID := uuid.NewString()
	payment.TransactionUuid = newUUID
	r.data = append(r.data, payment)
	log.Printf("Оплата прошла успешно, transaction_uuid: %s\n", payment.TransactionUuid)

	return payment.TransactionUuid, nil
}
