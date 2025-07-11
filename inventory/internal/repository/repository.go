package repository

import (
	"context"
	"github.com/dfg007star/go_rocket/inventory/internal/model"
	repoModel "github.com/dfg007star/go_rocket/inventory/internal/repository/model"
)

type PartRepository interface {
	Get(ctx context.Context, part repoModel.Part) (string, error)
	List(ctx context.Context, fitter repoModel.PartsFilter) ([]model.Part, error)
}
