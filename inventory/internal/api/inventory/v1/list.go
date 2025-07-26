package v1

import (
	"context"

	"github.com/dfg007star/go_rocket/inventory/internal/converter"
	"github.com/dfg007star/go_rocket/inventory/internal/model"
	inventoryV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	filter := &model.PartsFilter{}

	if req.Filter != nil {
		filter = &model.PartsFilter{
			Uuids:                 req.Filter.Uuids,
			Names:                 req.Filter.Names,
			Categories:            converter.ConvertCategories(req.Filter.Categories),
			ManufacturerCountries: req.Filter.ManufacturerCountries,
			Tags:                  req.Filter.Tags,
		}
	}

	parts, err := a.inventoryService.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	return converter.PartsModelToGrpcResponse(parts), nil
}
