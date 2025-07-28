package order

import (
	"fmt"
	"github.com/dfg007star/go_rocket/order/internal/config"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"

	"github.com/dfg007star/go_rocket/order/internal/migrator"
	def "github.com/dfg007star/go_rocket/order/internal/repository"
)

var _ def.OrderRepository = (*repository)(nil)

type repository struct {
	mu   sync.RWMutex
	data *pgx.Conn
}

func NewRepository(clientPostgres *pgx.Conn) *repository {
	fmt.Println(config.AppConfig().Postgres.MigrationDirectory())
	migratorRunner := migrator.NewMigrator(stdlib.OpenDB(*clientPostgres.Config().Copy()), config.AppConfig().Postgres.MigrationDirectory())

	err := migratorRunner.Up()
	if err != nil {
		panic(fmt.Errorf("error while migration db: %w", err))
	}

	return &repository{
		data: clientPostgres,
	}
}
