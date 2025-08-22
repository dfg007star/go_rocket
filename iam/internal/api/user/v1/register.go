package v1

import (
	"context"
	"time"

	"github.com/dfg007star/go_rocket/iam/internal/model"
	userV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/user/v1"
)

func (a *api) Register(ctx context.Context, req *userV1.RegisterRequest) (*userV1.RegisterResponse, error) {
	user := model.User{
		Login:               req.Login,
		Email:               req.Email,
		Password:            req.Password,
		NotificationMethods: ???,
		CreatedAt:           time.Now(),
	}
}
