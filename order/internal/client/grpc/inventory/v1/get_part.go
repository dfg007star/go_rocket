package v1

import (
	"context"

	clientConverter "github.com/dfg007star/go_rocket/order/internal/client/converter"
	"github.com/dfg007star/go_rocket/order/internal/model"
	grpcAuth "github.com/dfg007star/go_rocket/platform/pkg/middleware/grpc"
	generatedInventoryV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/inventory/v1"
)

func (c *client) GetPart(ctx context.Context, uuid string) (*model.Part, error) {
	// добавляем session-uuid в gRPC metadata для передачи в Inventory сервис
	ctx = grpcAuth.ForwardSessionUUIDToGRPC(ctx)

	part, err := c.generatedClient.GetPart(ctx, &generatedInventoryV1.GetPartRequest{Uuid: uuid})
	if err != nil {
		return nil, err
	}

	return clientConverter.PartToModel(part.Part), nil
}
