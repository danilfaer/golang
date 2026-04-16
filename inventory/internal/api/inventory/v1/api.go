package v1

import (
	"github.com/danilfaer/golang/inventory/internal/service"
	inventoryV1 "github.com/danilfaer/golang/shared/pkg/proto/inventory/v1"
)

type api struct {
	inventoryV1.UnimplementedInventoryServiceServer

	inventoryService service.InventoryService
}

func NewAPI(inventoryService service.InventoryService) *api {
	return &api{inventoryService: inventoryService}
}
