package v1

import (
	def "github.com/dfg007star/go_rocket/inventory/internal/client/grpc"
	generatedIAMV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/auth/v1"
)

var _ def.IAMClient = (*client)(nil)

type client struct {
	generatedClient generatedIAMV1.AuthServiceClient
}

func NewClient(generatedClient generatedIAMV1.AuthServiceClient) *client {
	return &client{generatedClient: generatedClient}
}
