package repository

import (
	"context"

	"github.com/danilfaer/golang/order/internal/repository/model"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *model.Order) (string, error)
	GetOrderByUuid(ctx context.Context, uuid string) (*model.Order, error)
	UpdateOrder(ctx context.Context, order *model.Order) error
}
