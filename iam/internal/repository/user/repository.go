package user

import (
	"github.com/jackc/pgx/v5"

	def "github.com/dfg007star/go_rocket/iam/internal/repository"
)

var _ def.UserRepository = (*repository)(nil)

type repository struct {
	data *pgx.Conn
}

func NewUserRepository(clientPostgres *pgx.Conn) *repository {
	return &repository{
		data: clientPostgres,
	}
}
