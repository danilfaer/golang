package service

import (
	"context"

	"github.com/danilfaer/golang/payment/internal/model"
)

type PaymentService interface {
	PayOrder(ctx context.Context, req model.Pay) (string, error)
}
