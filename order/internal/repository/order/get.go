package order

import (
	"context"

	orderModel "github.com/danilfaer/golang/order/internal/model"
	repoModel "github.com/danilfaer/golang/order/internal/repository/model"
)

func (r *repository) GetOrderByUuid(ctx context.Context, uuid string) (*repoModel.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, exists := r.orders[uuid]
	if !exists {
		return nil, orderModel.ErrOrderNotFound
	}

	return order, nil
}
