package part

import (
	"errors"

	"github.com/dfg007star/go_rocket/inventory/internal/model"
)

func (s *service) validatePart(part *model.Part) error {
	if part.Name == "" {
		return errors.New("part name is required")
	}

	if part.Price < 0 {
		return errors.New("price cannot be negative")
	}

	if part.StockQuantity < 0 {
		return errors.New("stock quantity cannot be negative")
	}

	if part.Category == model.UNKNOWN {
		return errors.New("category must be specified")
	}

	if part.Dimensions.Weight <= 0 {
		return errors.New("weight must be positive")
	}

	if part.Manufacturer.Name == "" {
		return errors.New("manufacturer name is required")
	}

	return nil
}
