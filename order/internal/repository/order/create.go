package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/danilfaer/golang/order/internal/repository/model"
)


func (r *repository) CreateOrder(ctx context.Context, req *model.Order) (string, error) {

	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return "", err
	}
	defer conn.Release()

	// Транзакция: единая точка фиксации, чтобы при добавлении связанных шагов не плодить частично записанные заказы.
	tx, err := conn.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("начать транзакцию: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	orderUUID := uuid.New().String()

	const q = `
		INSERT INTO orders (order_uuid, user_uuid, part_uuids, total_price, payment_method, status)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	
    _, err = tx.Exec(ctx, q, orderUUID, req.UserUUID, req.PartUuids, req.TotalPrice, req.PaymentMethod, req.Status)
	
	if err != nil {
		return "", fmt.Errorf("вставка заказа: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return "", fmt.Errorf("зафиксировать транзакцию: %w", err)
	}
	return orderUUID, nil
}
