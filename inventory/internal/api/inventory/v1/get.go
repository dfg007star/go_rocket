package v1

import (
	"context"
	"github.com/dfg007star/go_rocket/inventory/internal/converter"
	"github.com/dfg007star/go_rocket/inventory/internal/model"
	inventoryV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/inventory/v1"
)

func (a *api) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	part, err := a.inventoryService.Get(ctx, req.Uuid)

	if err != nil {
		return nil, model.ErrPartNotFound
	}

	return converter.PartModelToGrpcResponse(&part), nil
}
