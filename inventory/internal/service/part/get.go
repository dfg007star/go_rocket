package part

import (
	"context"

	"github.com/dfg007star/go_rocket/inventory/internal/model"
)

func (s *service) Get(ctx context.Context, uuid string) (*model.Part, error) {
	if uuid == "" {
		return nil, model.ErrUuidIsEmpty
	}

	part, err := s.inventoryRepository.Get(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return part, nil
}
