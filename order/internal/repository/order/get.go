package order

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	repoModel "github.com/danilfaer/golang/order/internal/repository/model"
)



func (r *repository) GetOrderByUuid(ctx context.Context, id string) (*repoModel.Order, error) {

	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	const q = `
		SELECT order_uuid, user_uuid, part_uuids, total_price, transaction_uuid, payment_method, status
		FROM orders
		WHERE order_uuid = $1
	`

	var order repoModel.Order
	
	err = conn.QueryRow(ctx,q, id).Scan(&order.OrderUUID, &order.UserUUID, &order.PartUuids, &order.TotalPrice, &order.TransactionUUID, &order.PaymentMethod, &order.Status)
	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}
		return nil, err
	}
	return &order, nil
}
