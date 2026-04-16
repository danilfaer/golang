package part

import (
	"context"

	"github.com/danilfaer/golang/inventory/internal/converter"
	"github.com/danilfaer/golang/inventory/internal/model"
)

func (s *service) ListParts(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error) {
	repoFilter := converter.ConvertModelPartsFilterToRepoPartsFilter(filter)
	parts, err := s.inventoryRepository.ListParts(ctx, repoFilter)
	if err != nil {
		return nil, err
	}

	return converter.ConvertRepoPartsToModelParts(parts), nil
}
