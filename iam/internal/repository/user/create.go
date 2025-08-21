package user

import (
	"context"

	"github.com/dfg007star/go_rocket/iam/internal/model"
)

func (r *repository) Create(ctx context.Context, user *model.User) (string, error) {
	repoUser := converter
}
