package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/dfg007star/go_rocket/iam/internal/config/env"
)

var appConfig *config

type config struct {
	Logger   LoggerConfig
	IamGRPC  IamGRPCConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	Session  SessionConfig
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

	iamGRPCCfg, err := env.NewIamGRPCConfig()
	if err != nil {
		return err
	}

	postgresCfg, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	redisCfg, err := env.NewRedisConfig()
	if err != nil {
		return err
	}

	sessionCfg, err := env.NewSessionConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:   loggerCfg,
		IamGRPC:  iamGRPCCfg,
		Postgres: postgresCfg,
		Redis:    redisCfg,
		Session:  sessionCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
