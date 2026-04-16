package v1

import (
	"context"

	"github.com/danilfaer/golang/order/internal/client/converter"
	"github.com/danilfaer/golang/order/internal/model"
	genaratedInventoryV1 "github.com/danilfaer/golang/shared/pkg/proto/inventory/v1"
)

func (c *client) ListParts(ctx context.Context, filter model.PartsFilter) ([]*model.Part, error) {
	parts, err := c.generatedClient.ListParts(ctx, &genaratedInventoryV1.ListPartsRequest{
		Filter: converter.PartsFilterToProto(filter),
	})
	if err != nil {
		return nil, err
	}
	return converter.PartListProtoToModel(parts.Parts), nil
}
