package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/dfg007star/go_rocket/order/internal/model"
	"github.com/dfg007star/go_rocket/order/internal/repository/converter"
	repoModel "github.com/dfg007star/go_rocket/order/internal/repository/model"
)

func (r *repository) Create(ctx context.Context, userUuid string, parts []*model.Part) (*model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	partUuids := make([]string, 0, len(parts))
	totalPrice := 0.0

	for _, part := range parts {
		partUuids = append(partUuids, part.Uuid)
		totalPrice += part.Price
	}

	order := &repoModel.Order{
		OrderUuid:       uuid.New().String(),
		UserUuid:        userUuid,
		PartUuids:       partUuids,
		TotalPrice:      float32(totalPrice),
		TransactionUuid: nil,
		PaymentMethod:   nil,
		Status:          repoModel.PENDING_PAYMENT,
	}

	_, err := r.data.Exec(ctx, `
		INSERT INTO orders (
			order_uuid, 
			user_uuid, 
			part_uuids, 
			total_price, 
			status,
			created_at
		) VALUES ($1, $2, $3, $4, $5, $6)`,
		order.OrderUuid,
		order.UserUuid,
		order.PartUuids,
		order.TotalPrice,
		order.Status.String(),
		order.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return converter.RepoModelToOrder(order), nil
}
