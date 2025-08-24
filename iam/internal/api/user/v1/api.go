package v1

import (
	"github.com/dfg007star/go_rocket/iam/internal/service"
	userV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/user/v1"
)

type api struct {
	userV1.UnimplementedUserServiceServer

	userService service.UserService
}

func NewUserAPI(userService service.UserService) *api {
	return &api{
		userService: userService,
	}
}
