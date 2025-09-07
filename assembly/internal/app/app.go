package app

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/dfg007star/go_rocket/assembly/internal/config"
	assemblyMetrics "github.com/dfg007star/go_rocket/assembly/internal/metrics"
	"github.com/dfg007star/go_rocket/platform/pkg/closer"
	"github.com/dfg007star/go_rocket/platform/pkg/logger"
	loggerConfig "github.com/dfg007star/go_rocket/platform/pkg/logger"
	"github.com/dfg007star/go_rocket/platform/pkg/metrics"
)

type App struct {
	diContainer *diContainer
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
	errCh := make(chan error, 1)

	// Контекст для остановки всех горутин
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Консьюмер
	go func() {
		if err := a.runConsumer(ctx); err != nil {
			errCh <- errors.Errorf("consumer crashed: %v", err)
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
		a.initLogger,
		a.initMetrics,
		a.initCloser,
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

func (a *App) initMetrics(ctx context.Context) error {
	closer.AddNamed("Metrics Assembly", func(ctx context.Context) error {
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

	err = assemblyMetrics.InitMetrics()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) runConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 Assembly Kafka consumer running")

	err := a.diContainer.AssemblyConsumerService().RunConsumer(ctx)
	if err != nil {
		return err
	}

	return nil
}
