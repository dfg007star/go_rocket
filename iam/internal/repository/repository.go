package repository

import (
	"context"
	"time"

	"github.com/dfg007star/go_rocket/iam/internal/model"
)

type UserRepository interface {
	GetByUserUuid(ctx context.Context, userUuid *string) (*model.User, error)
	GetByUserLogin(ctx context.Context, login *string) (*model.User, error)
	Register(ctx context.Context, user *model.User) (*string, error)
}

type SessionRepository interface {
	Get(ctx context.Context, sessionUuid *string) (*model.User, error)
	Create(ctx context.Context, sessionUuid *string, user *model.User, ttl time.Duration) error
}
