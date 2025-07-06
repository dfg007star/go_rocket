package main

import (
	orderv1 "github.com/dfg007star/go_rocket/shared/pkg/openapi/order/v1"
	"sync"
	"time"
)

const (
	httpPort = "8080"
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

// OrderService provides thread-safe storage and management of orders
type OrderService struct {
	mu     sync.RWMutex
	orders map[string]*orderv1.OrderDto
}

// NewOrderService creates and returns a new initialized OrderService instance
func NewOrderService() *OrderService {
	return &OrderService{
		orders: make(map[string]*orderv1.OrderDto),
	}
}

func (s *OrderService) OrderByUuid(order_uuid string) *orderv1.OrderDto {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[order_uuid]
	if !ok {
		return nil
	}

	return order
}

func (s *OrderService) CreateOrder(order_uuid string) *orderv1.OrderDto {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[order_uuid]
	if !ok {
		return nil
	}

	return order
}

func (s *OrderService) PayOrder(order_uuid string) *orderv1.OrderDto {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[order_uuid]
	if !ok {
		return nil
	}

	return order
}

func (s *OrderService) CancelOrderByUuid(order_uuid string) *orderv1.OrderDto {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[order_uuid]
	if !ok {
		return nil
	}

	return order
}

// OrderHandler provides HTTP handlers for order operations
type OrderHandler struct {
	service *OrderService
}

// NewOrderHandler creates a new OrderHandler instance with the given OrderService
func NewOrderHandler(service *OrderService) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

//func (h *OrderHandler) OrderByUuid(_ context.Context, params orderv1.OrderByUuidParams) (orderv1.OrderByUuidRes, error) {
//	order := h.service.OrderByUuid(params.OrderUUID)
//	if order == nil {
//		return &orderv1.NotFoundError{
//			Code:    404,
//			Message: "Order with UUID <'" + params.OrderUUID + "'> not found",
//		}, nil
//	}
//
//	return order, nil
//}

//func (h *OrderHandler) CreateOrder(_ context.Context, req *order_v1.CreateOrderRequest) (order_v1.CreateOrderResponse, error) {}

func main() {
	//service := NewOrderService()

	//orderHandler := NewOrderHandler(service)

	// need all interface 5 :)
	//orderServer, err := order_v1.NewServer(orderHandler)
}
