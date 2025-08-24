package app

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/dfg007star/go_rocket/iam/internal/config"
	"github.com/dfg007star/go_rocket/platform/pkg/closer"
	"github.com/dfg007star/go_rocket/platform/pkg/grpc/health"
	"github.com/dfg007star/go_rocket/platform/pkg/logger"
	pgMigrator "github.com/dfg007star/go_rocket/platform/pkg/migrator/pg"
	authV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/auth/v1"
	userV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/user/v1"
)

type App struct {
	diContainer *diContainer
	grpcServer  *grpc.Server
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
	return a.runGRPCServer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
		a.initMigrator,
		a.initListener,
		a.initGRPCServer,
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
	listener, err := net.Listen("tcp", config.AppConfig().IamGRPC.Address())
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

func (a *App) initMigrator(ctx context.Context) error {
	migrator := pgMigrator.New(
		stdlib.OpenDB(*a.diContainer.PostgresClient(ctx).Config().Copy()),
		config.AppConfig().Postgres.MigrationDirectory(),
	)
	err := migrator.Up()
	if err != nil {
		panic(fmt.Errorf("failed to migrate db: %w", err))
	}

	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	closer.AddNamed("gRPC server", func(ctx context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	reflection.Register(a.grpcServer)

	// Регистрируем health service для проверки работоспособности
	health.RegisterService(a.grpcServer)

	userV1.RegisterUserServiceServer(a.grpcServer, a.diContainer.UserV1API(ctx))
	authV1.RegisterAuthServiceServer(a.grpcServer, a.diContainer.AuthV1API(ctx))

	return nil
}

func (a *App) runGRPCServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("🚀 gRPC IAMService server listening on %s", config.AppConfig().IamGRPC.Address()))

	err := a.grpcServer.Serve(a.listener)
	if err != nil {
		return err
	}

	return nil
}
