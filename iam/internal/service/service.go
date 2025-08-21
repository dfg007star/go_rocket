package service

import (
	"context"
	"github.com/dfg007star/go_rocket/iam/internal/model"
)

type UserService interface {
	Get(ctx context.Context, userUuid *string) (*model.User, error)
	Register(ctx context.Context, user *model.User) (*string, error)
}
