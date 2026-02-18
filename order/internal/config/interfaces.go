package config

import ()

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type OrderHTTPConfig interface {
	Address() string
	ReadHeaderTimeout() string
}

type PostgresConfig interface {
	URI() string
	DatabaseName() string
	MigrationDir() string
}

type InventoryGRPCConfig interface {
	Address() string
}

type PaymentGRPCConfig interface {
	Address() string
}
