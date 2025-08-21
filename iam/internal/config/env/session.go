package env

import (
	"github.com/caarlos0/env/v11"
)

type sessionEnvConfig struct {
	SessionTtl string `env:"SESSION_TTL,required"`
}

type sessionConfig struct {
	raw sessionEnvConfig
}

func NewSessionConfig() (*sessionConfig, error) {
	var raw sessionEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &sessionConfig{raw: raw}, nil
}

func (cfg *sessionConfig) SessionTtl() string {
	return cfg.raw.SessionTtl
}
