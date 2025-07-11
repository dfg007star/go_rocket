package part

import (
	def "github.com/dfg007star/go_rocket/inventory/internal/repository"
	repoModel "github.com/dfg007star/go_rocket/inventory/internal/repository/model"
	"sync"
)

var _ def.PartRepository = (*repository)(nil)

type repository struct {
	mu   sync.RWMutex
	data []repoModel.Part
}

func NewRepository() *repository {
	return &repository{
		data: []repoModel.Part{},
	}
}
