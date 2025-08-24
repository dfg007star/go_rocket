package grpc

import (
	"context"

	authV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/auth/v1"
)

type IAMClient interface {
	WhoAmI(ctx context.Context, req *authV1.WhoAmIRequest) (*authV1.WhoAmIResponse, error)
}
