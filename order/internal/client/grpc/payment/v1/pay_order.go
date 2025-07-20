package v1

import (
	"context"
	"fmt"

	"github.com/dfg007star/go_rocket/order/internal/client/converter"
	"github.com/dfg007star/go_rocket/order/internal/model"
	generatedPaymentV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/payment/v1"
)

func (c *client) PayOrder(ctx context.Context, paymentMethod model.PaymentMethod, orderUuid, userUuid string) (string, error) {
	res, err := c.generatedClient.PayOrder(ctx, &generatedPaymentV1.PayOrderRequest{
		OrderUuid:     orderUuid,
		UserUuid:      userUuid,
		PaymentMethod: converter.PaymentMethodToProto(paymentMethod),
	})
	fmt.Println("pay order client", err)
	if err != nil {
		return "", err
	}

	return res.TransactionUuid, nil
}
