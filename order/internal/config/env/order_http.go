package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type orderHTTPEnvConfig struct {
	Host string `env:"HTTP_HOST,required"`
	Port string `env:"HTTP_PORT,required"`
}

type orderHTTPConfig struct {
	raw orderHTTPEnvConfig
}

func NewOrderHTTPConfig () (*orderHTTPConfig, error) {
	var raw orderHTTPEnvConfig
	err := env.Parse(&raw)
	if err != nil {
		return nil, err
	}

	return &orderHTTPConfig{
		raw: raw,
	}, nil
}

func (cfg *orderHTTPConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}


