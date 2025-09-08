package order

import (
	"context"
	"strings"

	orderMetrics "github.com/dfg007star/go_rocket/order/internal/metrics"
	"github.com/dfg007star/go_rocket/platform/pkg/tracing"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/dfg007star/go_rocket/order/internal/model"
)

func (s *service) Pay(ctx context.Context, orderUuid string, method *model.PaymentMethod) (*model.Order, error) {
	ctx, getOrderSpan := tracing.StartSpan(ctx, "order.call_get_order",
		trace.WithAttributes(
			attribute.String("order.uuid", orderUuid),
		),
	)
	defer getOrderSpan.End()

	order, err := s.orderRepository.Get(ctx, orderUuid)
	if err != nil {
		getOrderSpan.RecordError(err)
		getOrderSpan.End()
		return nil, err
	}

	if order.Status == model.PAID {
		return nil, model.ErrOrderAlreadyPaid
	}

	if order.Status == model.CANCELLED {
		return nil, model.ErrOrderAlreadyCancelled
	}

	paymentMethod := method
	if order.PaymentMethod != nil {
		paymentMethod = order.PaymentMethod
	}

	// Создаем спан для вызова Payment сервиса
	ctx, orderPaySpan := tracing.StartSpan(ctx, "order.call_pay_order",
		trace.WithAttributes(
			attribute.String("order.uuid", order.OrderUuid),
			attribute.String("order.user_uuid", order.UserUuid),
			attribute.String("order.part_uuids", strings.Join(order.PartUuids, ",")),
		),
	)
	defer orderPaySpan.End()
	transactionUuid, err := s.paymentClient.PayOrder(ctx, paymentMethod, order.OrderUuid, order.UserUuid)
	if err != nil {
		orderPaySpan.RecordError(err)
		orderPaySpan.End()
		return nil, err
	}

	paidStatus := model.PAID
	updatedOrder, err := s.orderRepository.Update(ctx, &model.OrderUpdate{
		OrderUuid:       order.OrderUuid,
		TransactionUuid: &transactionUuid,
		Status:          &paidStatus,
		PaymentMethod:   paymentMethod,
	})
	if err != nil {
		return nil, err
	}

	var paymentMethodStr string
	if paymentMethod != nil {
		paymentMethodStr = paymentMethod.String()
	}

	event := model.OrderPaidEvent{
		EventUuid:       uuid.New().String(),
		OrderUuid:       order.OrderUuid,
		UserUuid:        order.UserUuid,
		PaymentMethod:   paymentMethodStr,
		TransactionUuid: transactionUuid,
	}
	err = s.orderProducerService.OrderPaid(ctx, event)
	if err != nil {
		return nil, err
	}

	// Бизнес-метрика: суммарная выручка
	orderMetrics.OrdersRevenueTotal.Add(ctx, float64(order.TotalPrice))

	return updatedOrder, nil
}
