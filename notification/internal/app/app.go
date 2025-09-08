package app

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/dfg007star/go_rocket/notification/internal/config"
	"github.com/dfg007star/go_rocket/platform/pkg/closer"
	"github.com/dfg007star/go_rocket/platform/pkg/logger"
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
	errCh := make(chan error, 2)

	// Контекст для остановки всех горутин
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Paid Консьюмер
	go func() {
		if err := a.runPaidConsumer(ctx); err != nil {
			errCh <- errors.Errorf("consumer crashed: %v", err)
		}
	}()

	// Assembled Консьюмер
	go func() {
		if err := a.runAssembledConsumer(ctx); err != nil {
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
		a.initCloser,
		a.initTelegramBot,
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

func (a *App) initLogger(ctx context.Context) error {
	conf := &logger.LoggerConf{
		LevelStr:           config.AppConfig().Logger.Level(),
		AsJSON:             config.AppConfig().Logger.AsJson(),
		EnableOTLP:         config.AppConfig().Logger.EnableOTLP(),
		OTLPEndpoint:       config.AppConfig().Logger.OTLPEndpoint(),
		ServiceName:        config.AppConfig().Logger.ServiceName(),
		ServiceEnvironment: config.AppConfig().Logger.ServiceEnvironment(),
	}

	return logger.Init(ctx, conf)
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initTelegramBot(ctx context.Context) error {
	// Получаем бота из DI контейнера
	telegramBot := a.diContainer.TelegramBot(ctx)

	// Регистрируем обработчик для активации бота
	telegramBot.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, func(ctx context.Context, b *bot.Bot, update *models.Update) {
		logger.Info(ctx, "chat id", zap.Int64("chat_id", update.Message.Chat.ID))

		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "🚀 Notification Bot активирован! Теперь вы будете получать уведомления о новых действиях на площадке.",
		})
		if err != nil {
			logger.Error(ctx, "Failed to send activation message", zap.Error(err))
		}
	})

	// Запускаем бота в фоне
	go func() {
		logger.Info(ctx, "🤖 Telegram bot started...")
		telegramBot.Start(ctx)
	}()

	return nil
}

func (a *App) runAssembledConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 Assembly Kafka consumer running")

	err := a.diContainer.OrderAssembledConsumerService(ctx).RunConsumer(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runPaidConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 Paid Kafka consumer running")

	err := a.diContainer.OrderPaidConsumerService(ctx).RunConsumer(ctx)
	if err != nil {
		return err
	}

	return nil
}
