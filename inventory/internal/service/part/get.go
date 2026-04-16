package part

import (
	"context"

	"github.com/danilfaer/golang/inventory/internal/converter"
	"github.com/danilfaer/golang/inventory/internal/model"
)

func (s *service) GetPart(ctx context.Context, uuid string) (*model.Part, error) {
	part, err := s.inventoryRepository.GetPart(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return converter.ConvertRepoPartToModelPart(part), nil
}
