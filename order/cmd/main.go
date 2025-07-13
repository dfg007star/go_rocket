package main

import (
	"context"
	"errors"
	"fmt"
	orderAPI "github.com/dfg007star/go_rocket/order/internal/api/order/v1"
	inventoryServiceClient "github.com/dfg007star/go_rocket/order/internal/client/grpc/inventory/v1"
	paymentServiceClient "github.com/dfg007star/go_rocket/order/internal/client/grpc/payment/v1"
	orderRepository "github.com/dfg007star/go_rocket/order/internal/repository/order"
	orderService "github.com/dfg007star/go_rocket/order/internal/service/order"
	orderV1 "github.com/dfg007star/go_rocket/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/payment/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	httpPort = "8080"
	// Таймауты для HTTP-сервера
	readHeaderTimeout    = 5 * time.Second
	shutdownTimeout      = 10 * time.Second
	grpcInventoryAddress = "localhost:50051"
	grpcPaymentAddress   = "localhost:50052"
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

func (s *OrderService) OrderUpdate(orderUuid string, order *orderV1.OrderDto) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	s.orders[orderUuid] = order
}

func (s *OrderService) CreateOrder(userUuid string, parts []*inventoryV1.Part) *orderV1.OrderDto {
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
	inventoryService inventoryV1.InventoryServiceClient
	paymentService   paymentV1.PaymentServiceClient
}

// NewOrderHandler creates a new OrderHandler instance with the given OrderService
func NewOrderHandler(service *OrderService, grpcInventoryConn *grpc.ClientConn, grpcPaymentConn *grpc.ClientConn) *OrderHandler {
	return &OrderHandler{
		service:          service,
		inventoryService: inventoryV1.NewInventoryServiceClient(grpcInventoryConn),
		paymentService:   paymentV1.NewPaymentServiceClient(grpcPaymentConn),
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
	partsReq := inventoryV1.ListPartsRequest{Filter: &inventoryV1.PartsFilter{Uuids: req.PartUuids}}
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

func (h *OrderHandler) PayOrder(_ context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	order := h.service.OrderByUuid(params.OrderUUID)
	if order == nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Order with UUID <'" + params.OrderUUID + "'> not found",
		}, nil
	}

	l, err := convertPaymentMethod(req.PaymentMethod)
	if err != nil {
		return &orderV1.NotFoundError{
			Code:    400,
			Message: "Payment method not found",
		}, nil
	}

	paymentReq := paymentV1.PayOrderRequest{UserUuid: order.UserUUID, OrderUuid: order.OrderUUID, PaymentMethod: l}
	payment, err := h.paymentService.PayOrder(context.Background(), &paymentReq)
	if err != nil {
		return &orderV1.NotFoundError{
			Code:    400,
			Message: "Payment processing failed: " + err.Error(),
		}, nil
	}

	order.SetTransactionUUID(orderV1.OptString{
		Set:   true,
		Value: payment.TransactionUuid,
	})
	order.SetPaymentMethod(orderV1.OptOrderDtoPaymentMethod{
		Set:   true,
		Value: orderV1.OrderDtoPaymentMethod(req.PaymentMethod),
	})
	order.SetStatus(orderV1.OrderDtoStatusPAID)

	return &orderV1.PayOrderResponse{
		TransactionUUID: payment.TransactionUuid,
	}, nil
}

func (h *OrderHandler) CancelOrderByUuid(_ context.Context, params orderV1.CancelOrderByUuidParams) (orderV1.CancelOrderByUuidRes, error) {
	order := h.service.OrderByUuid(params.OrderUUID)
	if order == nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Order with UUID <'" + params.OrderUUID + "'> not found",
		}, nil
	}

	if order.Status == orderV1.OrderDtoStatusPAID {
		return &orderV1.ConflictError{
			Code:    409,
			Message: "Order with UUID <'" + params.OrderUUID + "'> already paid and cannot be cancelled",
		}, nil
	}

	order.SetStatus(orderV1.OrderDtoStatusCANCELLED)

	return &orderV1.CancelOrderByUuidNoContent{}, nil
}

func (h *OrderHandler) NewError(_ context.Context, err error) *orderV1.GenericErrorStatusCode {
	return &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderV1.GenericError{
			Code:    orderV1.NewOptInt(http.StatusInternalServerError),
			Message: orderV1.NewOptString(err.Error()),
		},
	}
}

func newInventoryClient() (inventoryV1.InventoryServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		grpcInventoryAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, err
	}

	client := inventoryV1.NewInventoryServiceClient(conn)
	return client, conn, nil
}

func newPaymentClient() (paymentV1.PaymentServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		grpcPaymentAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, err
	}

	client := paymentV1.NewPaymentServiceClient(conn)
	return client, conn, nil
}

func main() {
	inventoryClient, _, err := newInventoryClient()
	if err != nil {
		panic(err)
	}
	paymentClient, _, err := newPaymentClient()
	if err != nil {
		panic(err)
	}

	repo := orderRepository.NewRepository()
	service := orderService.NewOrderService(
		repo,
		inventoryServiceClient.NewClient(inventoryClient),
		paymentServiceClient.NewClient(paymentClient),
	)
	api := orderAPI.NewApi(service)

	orderServer, err := orderV1.NewServer(api)

	if err != nil {
		log.Fatalf("ошибка создания сервера OpenAPI: %v", err)
	}

	// Инициализируем роутер Chi
	r := chi.NewRouter()

	// Добавляем middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	// Монтируем обработчики OpenAPI
	r.Mount("/", orderServer)

	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}

func convertPaymentMethod(method orderV1.PayOrderRequestPaymentMethod) (paymentV1.PaymentMethod, error) {
	switch method {
	case orderV1.PayOrderRequestPaymentMethodPAYMENTMETHODUNSPECIFIED:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED, nil
	case orderV1.PayOrderRequestPaymentMethodPAYMENTMETHODCARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CARD, nil
	case orderV1.PayOrderRequestPaymentMethodPAYMENTMETHODSBP:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_SBP, nil
	case orderV1.PayOrderRequestPaymentMethodPAYMENTMETHODCREDITCARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD, nil
	case orderV1.PayOrderRequestPaymentMethodPAYMENTMETHODINVESTORMONEY:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY, nil
	default:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED,
			fmt.Errorf("unknown payment method: %v", method)

	}
}
