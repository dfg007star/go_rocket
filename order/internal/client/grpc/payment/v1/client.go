package v1

import (
	def "github.com/dfg007star/go_rocket/order/internal/client/grpc"
	generatedPaymentV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/payment/v1"
)

var _ def.PaymentClient = (*client)(nil)

type client struct {
	generatedClient generatedPaymentV1.PaymentServiceClient
}

func NewClient(generatedClient generatedPaymentV1.PaymentServiceClient) *client {
	return &client{generatedClient: generatedClient}
}
