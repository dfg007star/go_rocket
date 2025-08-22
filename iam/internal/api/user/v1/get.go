package v1

import (
	"context"

	userV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/user/v1"
)

func (a *api) GetUser(ctx context.Context, req *userV1.GetUserRequest) (*userV1.GetUserResponse, error) {
	userUuid := req.GetUserUuid()
	user, err := a.userService.Get(ctx, &userUuid)
	if err != nil {
		return nil, err
	}

	return *userV1.GetUserResponse{}
}
