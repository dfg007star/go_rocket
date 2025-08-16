package order

import (
	def "github.com/dfg007star/go_rocket/order/internal/repository"
	"github.com/jackc/pgx/v5"
)

var _ def.OrderRepository = (*repository)(nil)

type repository struct {
	data *pgx.Conn
}

func NewRepository(clientPostgres *pgx.Conn) *repository {
	return &repository{
		data: clientPostgres,
	}
}
