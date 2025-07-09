package v1

import (
	"github.com/dfg007star/go_rocket/payment/internal/service"
	payment_v1 "github.com/dfg007star/go_rocket/shared/pkg/proto/payment/v1"
)

type api struct {
	payment_v1.UnimplementedPaymentServiceServer

	paymentService service.PaymentService
}

func NewAPI(paymentService service.PaymentService) *api {
	return &api{
		paymentService: paymentService,
	}
}
