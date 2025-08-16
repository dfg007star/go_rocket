package config

import (
	"os"

	"github.com/dfg007star/go_rocket/order/internal/config/env"
	"github.com/joho/godotenv"
)

var appConfig *config

type config struct {
	Logger                 LoggerConfig
	InventoryGRPC          InventoryGRPCConfig
	PaymentGRPC            PaymentGRPCConfig
	OrderHTTP              OrderHTTPConfig
	Postgres               PostgresConfig
	Kafka                  KafkaConfig
	OrderAssembledConsumer OrderAssembledConsumerConfig
	OrderPaidProducer      OrderPaidProducerConfig
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

	inventoryGRPCCfg, err := env.NewInventoryGRPCConfig()
	if err != nil {
		return err
	}

	paymentGRPCCfg, err := env.NewPaymentGRPCConfig()
	if err != nil {
		return err
	}

	orderHTTPCfg, err := env.NewOrderHTTPConfig()
	if err != nil {
		return err
	}

	PostgresCfg, err := env.NewPostgresConfig()
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

	orderPaidProducerConfig, err := env.NewOrderPaidProducerConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:                 loggerCfg,
		InventoryGRPC:          inventoryGRPCCfg,
		PaymentGRPC:            paymentGRPCCfg,
		OrderHTTP:              orderHTTPCfg,
		Postgres:               PostgresCfg,
		Kafka:                  kafkaConfig,
		OrderAssembledConsumer: orderAssembledConsumerConfig,
		OrderPaidProducer:      orderPaidProducerConfig,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
