package v1

import "github.com/dfg007star/go_rocket/order/internal/service"

type api struct {
	orderService service.OrderService
}

func NewApi(orderService service.OrderService) *api {
	return &api{
		orderService: orderService,
	}
}
