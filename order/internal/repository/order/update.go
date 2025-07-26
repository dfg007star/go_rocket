package order

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	"github.com/dfg007star/go_rocket/order/internal/model"
	"github.com/dfg007star/go_rocket/order/internal/repository/converter"
)

func (r *repository) Update(ctx context.Context, orderUpdate *model.OrderUpdate) (*model.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	repoOrderUpdate := converter.OrderUpdateToRepoOrderUpdate(orderUpdate)

	queryBuilder := sq.Update("orders").
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"order_uuid": repoOrderUpdate.OrderUuid})

	if repoOrderUpdate.TransactionUuid != nil {
		queryBuilder = queryBuilder.Set("transaction_uuid", *repoOrderUpdate.TransactionUuid)
	}

	if repoOrderUpdate.PaymentMethod != nil {
		queryBuilder = queryBuilder.Set("payment_method", repoOrderUpdate.PaymentMethod.String())
	}

	if repoOrderUpdate.Status != nil {
		queryBuilder = queryBuilder.Set("status", repoOrderUpdate.Status.String())
	}

	sql, args, err := queryBuilder.
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build update query: %w", err)
	}

	_, err = r.data.Exec(ctx, sql, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to execute order update: %w", err)
	}

	return r.Get(ctx, orderUpdate.OrderUuid)
}
