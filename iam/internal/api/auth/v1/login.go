package v1

import (
	"context"

	authV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/auth/v1"
	commonV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/common/v1"
)

func (a *api) Login(ctx context.Context, req *authV1.LoginRequest) (*authV1.LoginResponse, error) {
	sessionUuid, err := a.authService.Login(ctx, &req.Login, &req.Password)
	if err != nil {
		return nil, err
	}

	return &authV1.LoginResponse{
		SessionUuid: &commonV1.SessionUuid{
			SessionUuid: *sessionUuid,
		},
	}, nil
}
