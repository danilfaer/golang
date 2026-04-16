package order

import (
	"sync"

	"github.com/danilfaer/golang/order/internal/repository/model"
)

type repository struct {
	mu     sync.RWMutex
	orders map[string]*model.Order
}

func NewRepository() *repository {
	return &repository{
		orders: make(map[string]*model.Order),
	}
}
