package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/dfg007star/go_rocket/order/internal/config"
	"github.com/dfg007star/go_rocket/platform/pkg/closer"
	"github.com/dfg007star/go_rocket/platform/pkg/logger"
	pgMigrator "github.com/dfg007star/go_rocket/platform/pkg/migrator/pg"
	"github.com/jackc/pgx/v5/stdlib"
	"net"
	"net/http"
)

type App struct {
	diContainer *diContainer
	httpServer  *http.Server
	listener    net.Listener
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	return a.runHTTPServer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
		a.initListener,
		a.initMigrator,
		a.initHTTPServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson(),
	)
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initListener(_ context.Context) error {
	listener, err := net.Listen("tcp", config.AppConfig().InventoryGRPC.Address())
	if err != nil {
		return err
	}
	closer.AddNamed("TCP listener", func(ctx context.Context) error {
		lerr := listener.Close()
		if lerr != nil && !errors.Is(lerr, net.ErrClosed) {
			return lerr
		}

		return nil
	})

	a.listener = listener

	return nil
}

func (a *App) initMigrator(_ context.Context) error {
	migrator := pgMigrator.New(
		stdlib.OpenDB(*a.diContainer.postgresClient.Config().Copy()),
		config.AppConfig().Postgres.MigrationDirectory(),
	)
	err := migrator.Up()
	if err != nil {
		panic(fmt.Errorf("failed to migrate db: %w", err))
	}

	return nil
}

func (a *App) initHTTPServer(_ context.Context) error {

}
