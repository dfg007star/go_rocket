package config

type LoggerConfig interface {
	Level() string
	AsJson() bool
	EnableOTLP() bool
	OTLPEndpoint() string
	ServiceName() string
	ServiceEnvironment() string
}

type InventoryGRPCConfig interface {
	Address() string
}

type MongoConfig interface {
	URI() string
	DatabaseName() string
}

type IamGRPCConfig interface {
	Address() string
}
