package part

import (
	"context"
	"github.com/dfg007star/go_rocket/inventory/internal/model"
	"time"
)

func (s *service) Create(ctx context.Context, part model.Part) (model.Part, error) {
	if err := s.validatePart(part); err != nil {
		return model.Part{}, err
	}

	part.CreatedAt = time.Now()

	return part, nil
}
