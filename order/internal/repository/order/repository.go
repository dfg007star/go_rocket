package order

import (
	"github.com/dfg007star/go_rocket/order/internal/model"
	def "github.com/dfg007star/go_rocket/order/internal/repository"
	"sync"
)

var _ def.OrderRepository = (*repository)(nil)

type repository struct {
	mu   sync.RWMutex
	data []model.Order
}

func NewRepository() *repository {
	return &repository{
		data: []model.Order{},
	}
}
