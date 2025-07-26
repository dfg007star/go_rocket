package part

import (
	"context"
	"time"

	"github.com/dfg007star/go_rocket/inventory/internal/model"
)

func (s *service) Create(ctx context.Context, part *model.Part) (*model.Part, error) {
	if err := s.validatePart(part); err != nil {
		return nil, err
	}

	part.CreatedAt = time.Now()

	return part, nil
}
