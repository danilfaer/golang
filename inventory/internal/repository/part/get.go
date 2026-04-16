package part

import (
	"context"

	"github.com/danilfaer/golang/inventory/internal/model"
	repoModel "github.com/danilfaer/golang/inventory/internal/repository/model"
)

func (r *repository) GetPart(_ context.Context, uuid string) (*repoModel.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	part, ok := r.parts[uuid]
	if !ok {
		return nil, model.ErrPartNotFound
	}
	return part, nil
}
