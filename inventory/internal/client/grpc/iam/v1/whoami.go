package v1

import (
	"context"

	authV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/auth/v1"
)

func (c *client) WhoAmI(ctx context.Context, req *authV1.WhoAmIRequest) (*authV1.WhoAmIResponse, error) {
	request, err := c.generatedClient.WhoAmI(ctx, req)
	if err != nil {
		return nil, err
	}

	return request, nil
}
