package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/dfg007star/go_rocket/order/internal/config"
	orderMetrics "github.com/dfg007star/go_rocket/order/internal/metrics"
	"github.com/dfg007star/go_rocket/platform/pkg/closer"
	"github.com/dfg007star/go_rocket/platform/pkg/logger"
	loggerConfig "github.com/dfg007star/go_rocket/platform/pkg/logger"
	"github.com/dfg007star/go_rocket/platform/pkg/metrics"
	middlewareHTTP "github.com/dfg007star/go_rocket/platform/pkg/middleware/http"
	pgMigrator "github.com/dfg007star/go_rocket/platform/pkg/migrator/pg"
)

type App struct {
	diContainer *diContainer
	httpServer  *http.Server
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
	// Канал для ошибок от компонентов
	errCh := make(chan error, 2)

	// Контекст для остановки всех горутин
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Консьюмер
	go func() {
		if err := a.runConsumer(ctx); err != nil {
			errCh <- errors.Errorf("consumer crashed: %v", err)
		}
	}()

	go func() {
		if err := a.runHTTPServer(ctx); err != nil {
			errCh <- errors.Errorf("http server crashed: %v", err)
		}
	}()

	// Ожидание либо ошибки, либо завершения контекста (например, сигнал SIGINT/SIGTERM)
	select {
	case <-ctx.Done():
		logger.Info(ctx, "Shutdown signal received")
	case err := <-errCh:
		logger.Error(ctx, "Component crashed, shutting down", zap.Error(err))
		// Триггерим cancel, чтобы остановить второй компонент
		cancel()
		// Дождись завершения всех задач (если есть graceful shutdown внутри)
		<-ctx.Done()
		return err
	}

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initMetrics,
		a.initLogger,
		a.initCloser,
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

func (a *App) initMetrics(ctx context.Context) error {
	closer.AddNamed("Metrics Order", func(ctx context.Context) error {
		err := metrics.Shutdown(ctx)
		if err != nil {
			return err
		}
		return nil
	})

	err := metrics.InitProvider(ctx, config.AppConfig().Metrics)
	if err != nil {
		return err
	}

	err = orderMetrics.InitMetrics()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initLogger(_ context.Context) error {
	conf := &loggerConfig.LoggerConf{
		LevelStr:           config.AppConfig().Logger.Level(),
		AsJSON:             config.AppConfig().Logger.AsJson(),
		EnableOTLP:         config.AppConfig().Logger.EnableOTLP(),
		OTLPEndpoint:       config.AppConfig().Logger.OTLPEndpoint(),
		ServiceName:        config.AppConfig().Logger.ServiceName(),
		ServiceEnvironment: config.AppConfig().Logger.ServiceEnvironment(),
	}

	return logger.Init(conf)
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
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

func (a *App) initHTTPServer(ctx context.Context) error {
	// Инициализируем роутер Chi
	r := chi.NewRouter()

	// Добавляем middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	sessionUuidMiddleware := middlewareHTTP.NewAuthMiddleware(a.diContainer.IamClient(ctx))
	r.Use(sessionUuidMiddleware.Handle)
	r.Use(middleware.Timeout(10 * time.Second))

	// Монтируем обработчики OpenAPI
	r.Mount("/", a.diContainer.OrderV1API(ctx))

	a.httpServer = &http.Server{
		Addr:              config.AppConfig().OrderHTTP.Address(),
		Handler:           r,
		ReadHeaderTimeout: config.AppConfig().OrderHTTP.ReadTimeout(),
	}

	closer.AddNamed("HTTP server", func(ctx context.Context) error {
		err := a.httpServer.Shutdown(ctx)
		if err != nil {
			return err
		}
		return nil
	})

	return nil
}

func (a *App) runHTTPServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("🚀 HTTP OrderService server listening on %s", config.AppConfig().OrderHTTP.Address()))

	err := a.httpServer.ListenAndServe()
	if err != nil {
		logger.Error(ctx, "failed to start http server")
		return err
	}

	return nil
}

func (a *App) runConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 Order Kafka consumer running")

	err := a.diContainer.OrderConsumerService(ctx).RunConsumer(ctx)
	if err != nil {
		return err
	}

	return nil
}
