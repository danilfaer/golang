package part

import "github.com/danilfaer/golang/inventory/internal/repository"

type service struct {
	inventoryRepository repository.InventoryRepository
}

func NewService(inventoryRepository repository.InventoryRepository) *service {
	return &service{inventoryRepository: inventoryRepository}
}
