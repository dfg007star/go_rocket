package payment

import (
	def "github.com/dfg007star/go_rocket/payment/internal/repository"
	repoModel "github.com/dfg007star/go_rocket/payment/internal/repository/model"
	"sync"
)

var _ def.PaymentRepository = (*repository)(nil)

type repository struct {
	mu   sync.RWMutex
	data []repoModel.Payment
}

func NewRepository() *repository {
	return &repository{
		data: []repoModel.Payment{},
	}
}
