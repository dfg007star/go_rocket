package v1

import (
	"context"

	"github.com/dfg007star/go_rocket/iam/internal/converter"
	authV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/auth/v1"
)

func (a *api) WhoAmI(ctx context.Context, req *authV1.WhoAmIRequest) (*authV1.WhoAmIResponse, error) {
	sessionUuid := req.SessionUuid.SessionUuid
	user, err := a.authService.WhoAmI(ctx, &sessionUuid)
	if err != nil {
		return nil, err
	}

	return converter.UserToWhoAmIResponse(user), nil
}
