// Code generated by ogen, DO NOT EDIT.

package order_v1

import (
	"context"
)

// Handler handles operations described by OpenAPI v3 specification.
type Handler interface {
	// CancelOrderByUuid implements CancelOrderByUuid operation.
	//
	// Cancel Order by UUID.
	//
	// POST /api/v1/orders/{order_uuid}/cancel
	CancelOrderByUuid(ctx context.Context, params CancelOrderByUuidParams) (CancelOrderByUuidRes, error)
	// CreateOrder implements CreateOrder operation.
	//
	// Create Order.
	//
	// POST /api/v1/orders
	CreateOrder(ctx context.Context, req *CreateOrderRequest) (CreateOrderRes, error)
	// OrderByUuid implements OrderByUuid operation.
	//
	// Get Order by UUID.
	//
	// GET /api/v1/orders/{order_uuid}
	OrderByUuid(ctx context.Context, params OrderByUuidParams) (OrderByUuidRes, error)
	// PayOrder implements PayOrder operation.
	//
	// Pay Order.
	//
	// POST /api/v1/orders/{order_uuid}/pay
	PayOrder(ctx context.Context, req *PayOrderRequest, params PayOrderParams) (PayOrderRes, error)
	// NewError creates *GenericErrorStatusCode from error returned by handler.
	//
	// Used for common default response.
	NewError(ctx context.Context, err error) *GenericErrorStatusCode
}

// Server implements http server based on OpenAPI v3 specification and
// calls Handler to handle requests.
type Server struct {
	h Handler
	baseServer
}

// NewServer creates new Server.
func NewServer(h Handler, opts ...ServerOption) (*Server, error) {
	s, err := newServerConfig(opts...).baseServer()
	if err != nil {
		return nil, err
	}
	return &Server{
		h:          h,
		baseServer: s,
	}, nil
}
