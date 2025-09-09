package env

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type metricEnvConfig struct {
	CollectorEndpoint string        `env:"METRIC_COLLECTOR_ENDPOINT,required"`
	CollectorInterval time.Duration `env:"METRIC_COLLECTOR_INTERVAL,required"`
}

type metricConfig struct {
	raw metricEnvConfig
}

func NewMetricConfig() (*metricConfig, error) {
	var raw metricEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &metricConfig{raw: raw}, nil
}

func (cfg *metricConfig) CollectorEndpoint() string {
	return cfg.raw.CollectorEndpoint
}

func (cfg *metricConfig) CollectorInterval() time.Duration {
	return cfg.raw.CollectorInterval
}
