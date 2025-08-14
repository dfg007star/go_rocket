package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/dfg007star/go_rocket/notification/internal/config/env"
)

var appConfig *config

type config struct {
	Logger                 LoggerConfig
	Kafka                  KafkaConfig
	OrderAssembledConsumer OrderAssembledConsumerConfig
	OrderPaidConsumer      OrderPaidConsumerConfig
	TelegramBot            TelegramBotConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	kafkaConfig, err := env.NewKafkaConfig()
	if err != nil {
		return err
	}

	orderAssembledConsumerConfig, err := env.NewOrderAssembledConsumerConfig()
	if err != nil {
		return err
	}

	orderPaidConsumerConfig, err := env.NewOrderPaidConsumerConfig()
	if err != nil {
		return err
	}

	telegramBotConfig, err := env.NewTelegramBotConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:                 loggerCfg,
		Kafka:                  kafkaConfig,
		OrderAssembledConsumer: orderAssembledConsumerConfig,
		OrderPaidConsumer:      orderPaidConsumerConfig,
		TelegramBot:            telegramBotConfig,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
