package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"

	"github.com/dfg007star/go_rocket/order/internal/model"
	"github.com/dfg007star/go_rocket/order/internal/repository/converter"
	repoModel "github.com/dfg007star/go_rocket/order/internal/repository/model"
)

func (r *repository) Get(ctx context.Context, orderUuid string) (*model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	query, args, err := squirrel.Select(
		"order_uuid",
		"user_uuid",
		"part_uuids",
		"total_price",
		"transaction_uuid",
		"payment_method",
		"status",
		"created_at",
		"updated_at",
	).
		From("orders").
		Where(squirrel.Eq{"order_uuid": orderUuid}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var dbOrder repoModel.Order
	var statusStr string
	err = r.data.QueryRow(ctx, query, args...).Scan(
		&dbOrder.OrderUuid,
		&dbOrder.UserUuid,
		&dbOrder.PartUuids,
		&dbOrder.TotalPrice,
		&dbOrder.TransactionUuid,
		&dbOrder.PaymentMethod,
		&statusStr,
		&dbOrder.CreatedAt,
		&dbOrder.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	dbOrder.Status = repoModel.StatusFromString(statusStr)

	return converter.RepoModelToOrder(&dbOrder), nil
}
