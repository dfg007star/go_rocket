package order

import (
	"context"

	"github.com/dfg007star/go_rocket/order/internal/model"
	"github.com/dfg007star/go_rocket/order/internal/repository/converter"
	repoModel "github.com/dfg007star/go_rocket/order/internal/repository/model"
	"github.com/google/uuid"
)

func (r *repository) Create(ctx context.Context, userUuid string, parts []model.Part) (model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	partUuids := make([]string, len(parts))
	totalPrice := 0.0

	for _, part := range parts {
		partUuids = append(partUuids, part.Uuid)
		totalPrice += part.Price
	}

	orderUuid := uuid.New().String()
	order := repoModel.Order{
		OrderUuid:       orderUuid,
		UserUuid:        userUuid,
		PartUuids:       partUuids,
		TotalPrice:      float32(totalPrice),
		TransactionUuid: nil,
		PaymentMethod:   nil,
		Status:          repoModel.PENDING_PAYMENT,
	}
	r.data = append(r.data, order)

	return converter.RepoModelToOrder(order), nil
}
