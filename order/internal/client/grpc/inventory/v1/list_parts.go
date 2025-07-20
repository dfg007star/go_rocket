package v1

import (
	"context"
	clientConverter "github.com/dfg007star/go_rocket/order/internal/client/converter"
	"github.com/dfg007star/go_rocket/order/internal/model"
	generatedInventoryV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/inventory/v1"
)

func (c *client) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	parts, err := c.generatedClient.ListParts(ctx, &generatedInventoryV1.ListPartsRequest{
		Filter: clientConverter.PartsFilterToProto(filter),
	})

	if err != nil {
		return nil, err
	}

	return clientConverter.PartsListToModel(parts.Parts), nil
}
