package order

import (
	"context"
	"fmt"

	"github.com/danilfaer/golang/order/internal/repository/model"
	"github.com/jackc/pgx/v5"
)

func (r *repository) UpdateOrder(ctx context.Context, order *model.Order) error {

	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	const q = `
		UPDATE orders 
		SET user_uuid = $1, part_uuids = $2, total_price = $3, transaction_uuid = $4, payment_method = $5, status = $6, updated_at = NOW()
		WHERE order_uuid = $7
	`

	tag, err := tx.Exec(ctx, q,
		order.UserUUID,
		order.PartUuids,
		order.TotalPrice,
		order.TransactionUUID,
		string(order.PaymentMethod),
		string(order.Status),
		order.OrderUUID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		// Сервис ожидает тот же сигнал «не найден», что и при чтении по UUID.
		return pgx.ErrNoRows
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("зафиксировать транзакцию: %w", err)
	}
	return nil
}
