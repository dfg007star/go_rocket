package part

import (
	"github.com/dfg007star/go_rocket/inventory/internal/repository"
	def "github.com/dfg007star/go_rocket/inventory/internal/service"
)

var _ def.InventoryService = (*service)(nil)

type service struct {
	inventoryRepository repository.PartRepository
}

func NewService(inventoryRepository repository.PartRepository) *service {
	return &service{
		inventoryRepository: inventoryRepository,
	}
}
