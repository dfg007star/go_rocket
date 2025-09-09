package config

import "time"

type LoggerConfig interface {
	Level() string
	AsJson() bool
	EnableOTLP() bool
	OTLPEndpoint() string
	ServiceName() string
	ServiceEnvironment() string
}

type IamGRPCConfig interface {
	Address() string
}

type PostgresConfig interface {
	URI() string
	MigrationDirectory() string
}

type RedisConfig interface {
	Address() string
	ConnectionTimeout() time.Duration
	MaxIdle() int
	IdleTimeout() time.Duration
	CacheTTL() time.Duration
}

type SessionConfig interface {
	SessionTtl() time.Duration
}
