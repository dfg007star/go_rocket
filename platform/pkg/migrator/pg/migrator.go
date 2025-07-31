package pg

import (
	"database/sql"
	"github.com/dfg007star/go_rocket/platform/pkg/migrator"
	"github.com/pressly/goose/v3"
)

type PostgresMigrator struct {
	migrator.Migrator
}

func New(db *sql.DB, migrationsDir string) *PostgresMigrator {
	return &PostgresMigrator{
		Migrator: migrator.Init(db, migrationsDir),
	}
}

func (m *PostgresMigrator) Up() error {
	return goose.Up(m.GetDB(), m.GetMigrationsDir())
}

func (m *PostgresMigrator) Down() error {
	return goose.Down(m.GetDB(), m.GetMigrationsDir())
}
