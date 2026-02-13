package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/mllbll/space-manufacture/order/internal/config/env"
)


var appConfig *config

type config struct {
	Logger LoggerConfig
	OrderHTTTP OrderHTTPConfig
	Postgres PostgresConfig
	InventoryGRPC InventoryGRPCConfig
	PaymentGRPC PaymentGRPCConfig
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

	orderHTTPCfg, err := env.NewOrderHTTPConfig()
	if err != nil {
		return err
	}

	postgresCfg, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	paymentCfg, err := env.NewPaymentGRPCConfig()
	if err != nil {
		return err
	}

	invetoryCfg, err := env.NewInventoryGRPCConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger: loggerCfg,
		OrderHTTTP: orderHTTPCfg,
		Postgres: postgresCfg,
		InventoryGRPC: invetoryCfg,
		PaymentGRPC: paymentCfg,
	}

	return nil

}

func AppConfig() *config {
	return appConfig
}
