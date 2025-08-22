package v1

import (
	"github.com/dfg007star/go_rocket/iam/internal/service"
	authV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/auth/v1"
)

type api struct {
	authV1.UnimplementedAuthServiceServer

	authService service.AuthService
}

func NewAuthAPI(authService service.AuthService) *api {
	return &api{
		authService: authService,
	}
}
