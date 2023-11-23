package config

import (
	"github.com/caarlos0/env/v10"
	"github.com/cockroachdb/errors"
	"github.com/zsmartex/pkg/v2/config"
	"go.uber.org/fx"
)

var Module = fx.Module("config",
	fx.Provide(
		NewConfig,
		provideConfig,
	),
)

type Config struct {
	HTTP     config.HTTP
	Database config.DB `envPrefix:"DATABASE_"`
}

type ConfigOut struct {
	fx.Out

	HTTP     config.HTTP `name:"http_server"`
	Database config.DB
}

func provideConfig(config *Config) ConfigOut {
	return ConfigOut{
		HTTP:     config.HTTP,
		Database: config.Database,
	}
}

func NewConfig() (*Config, error) {
	conf := new(Config)

	if err := env.Parse(conf); err != nil {
		return nil, errors.Newf("parse config: %v", err)
	}

	return conf, nil
}
