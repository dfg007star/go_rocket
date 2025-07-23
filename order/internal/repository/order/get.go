package order

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dfg007star/go_rocket/order/internal/repository/converter"
	repoModel "github.com/dfg007star/go_rocket/order/internal/repository/model"

	"github.com/dfg007star/go_rocket/order/internal/model"
)

func (r *repository) Get(ctx context.Context, orderUuid string) (*model.Order, error) {
	const query = `
		SELECT *
		FROM orders
		WHERE order_uuid = $1
	`

	var dbOrder repoModel.Order

	err := r.data.QueryRow(ctx, query, orderUuid).Scan(
		&dbOrder.OrderUuid,
		&dbOrder.UserUuid,
		&dbOrder.PartUuids,
		&dbOrder.TotalPrice,
		&dbOrder.TransactionUuid,
		&dbOrder.PaymentMethod,
		&dbOrder.Status,
		&dbOrder.CreatedAt,
		&dbOrder.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model.ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	fmt.Println("GET", dbOrder)

	return converter.RepoModelToOrder(&dbOrder), nil
}
