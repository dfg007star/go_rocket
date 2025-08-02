package order

import (
	"github.com/jackc/pgx/v5"

	def "github.com/dfg007star/go_rocket/order/internal/repository"
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
