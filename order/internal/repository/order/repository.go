package order

import (
	"sync"

	def "github.com/dfg007star/go_rocket/order/internal/repository"
	repoModel "github.com/dfg007star/go_rocket/order/internal/repository/model"
)

var _ def.OrderRepository = (*repository)(nil)

type repository struct {
	mu   sync.RWMutex
	data []repoModel.Order
}

func NewRepository() *repository {
	return &repository{
		data: []repoModel.Order{},
	}
}
