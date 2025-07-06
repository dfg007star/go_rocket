package main

import (
	"context"
	orderV1 "github.com/dfg007star/go_rocket/shared/pkg/openapi/order/v1"
	partV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/part/v1"
	paymentV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/payment/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"sync"
	"time"
)

const (
	httpPort = "8080"
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
	grpcAddress       = "localhost:50051"
)

// OrderService provides thread-safe storage and management of orders
type OrderService struct {
	mu     sync.RWMutex
	orders map[string]*orderV1.OrderDto
}

// NewOrderService creates and returns a new initialized OrderService instance
func NewOrderService() *OrderService {
	return &OrderService{
		orders: make(map[string]*orderV1.OrderDto),
	}
}

func (s *OrderService) OrderByUuid(orderUuid string) *orderV1.OrderDto {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[orderUuid]
	if !ok {
		return nil
	}

	return order
}

func (s *OrderService) CreateOrder(userUuid string, parts []*partV1.Part) *orderV1.OrderDto {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cp := make([]string, len(parts))
	totalPrice := 0.0

	for _, part := range parts {
		cp = append(cp, part.GetUuid())
		totalPrice += part.GetPrice()
	}
	orderUuid := uuid.New().String()
	order := &orderV1.OrderDto{
		OrderUUID:       orderUuid,
		UserUUID:        userUuid,
		PartUuids:       cp,
		TotalPrice:      float32(totalPrice),
		TransactionUUID: orderV1.OptString{Set: false},
		PaymentMethod:   orderV1.OptOrderDtoPaymentMethod{Set: false},
		Status:          orderV1.OrderDtoStatusPENDINGPAYMENT,
	}

	s.orders[orderUuid] = order
	return order
}

func (s *OrderService) PayOrder(order_uuid string) *orderV1.OrderDto {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[order_uuid]
	if !ok {
		return nil
	}

	return order
}

func (s *OrderService) CancelOrderByUuid(order_uuid string) *orderV1.OrderDto {
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
	service          *OrderService
	inventoryService partV1.PartServiceClient
	paymentService   paymentV1.PaymentServiceClient
}

// NewOrderHandler creates a new OrderHandler instance with the given OrderService
func NewOrderHandler(service *OrderService, grpcConn *grpc.ClientConn) *OrderHandler {
	return &OrderHandler{
		service:          service,
		inventoryService: partV1.NewPartServiceClient(grpcConn),
		paymentService:   paymentV1.NewPaymentServiceClient(grpcConn),
	}
}

func (h *OrderHandler) OrderByUuid(_ context.Context, params orderV1.OrderByUuidParams) (orderV1.OrderByUuidRes, error) {
	order := h.service.OrderByUuid(params.OrderUUID)
	if order == nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Order with UUID <'" + params.OrderUUID + "'> not found",
		}, nil
	}

	return order, nil
}

func (h *OrderHandler) CreateOrder(_ context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	partsReq := partV1.ListPartsRequest{Filter: &partV1.PartsFilter{Uuids: req.PartUuids}}
	parts, err := h.inventoryService.ListParts(context.Background(), &partsReq)
	if err != nil {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "Failed to get parts from inventoryService",
		}, nil
	}

	if len(parts.Parts) != len(req.PartUuids) {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "Not all parts were found in inventoryService",
		}, nil
	}

	order := h.service.CreateOrder(req.UserUUID, parts.Parts)
	if order == nil {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "Cannot create order in inventoryService",
		}, nil
	}

	return &orderV1.CreateOrderResponse{
		OrderUUID:  order.OrderUUID,
		TotalPrice: order.TotalPrice,
	}, nil
}

func (h *OrderHandler) PayOrder(_ context.Context, req *orderV1.PayOrderRequest, params *orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	order := h.service.OrderByUuid(params.OrderUUID)
	if order == nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Order with UUID <'" + params.OrderUUID + "'> not found",
		}, nil
	}
	paymentReq := paymentV1.PayOrderRequest{UserUuid: order.UserUUID, OrderUuid: order.OrderUUID, PaymentMethod: }
	payment := h.paymentService.PayOrder(context.Background())

}

func main() {
	ctx := context.Background()

	// grpcConn настройка gRPC клиента
	grpcConn, err := grpc.NewClient(
		grpcAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return
	}
	defer func() {
		if cerr := grpcConn.Close(); cerr != nil {
			log.Printf("failed to close connect: %v", cerr)
		}
	}()

	service := NewOrderService()
	orderHandler := NewOrderHandler(service, grpcConn)

	orderServer, err := orderV1.NewServer(orderHandler)
}
