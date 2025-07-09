package payment

import (
	"github.com/dfg007star/go_rocket/payment/internal/repository"
	def "github.com/dfg007star/go_rocket/payment/internal/service"
)

var _ def.PaymentService = (*service)(nil)

type service struct {
	paymentRepository repository.PaymentRepository
}

func NewService(paymentRepository repository.PaymentRepository) *service {
	return &service{
		paymentRepository: paymentRepository,
	}
}
