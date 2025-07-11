package v1

import part_v1 "github.com/dfg007star/go_rocket/shared/pkg/proto/inventory/v1"

type api struct {
	part_v1.UnimplementedInventoryServiceServer

	inventoryService service.InventoryService
}

func NewApi(inventoryService service.InventoryService) *api {
	return &api{
		inventoryService: inventoryService,
	}
}
