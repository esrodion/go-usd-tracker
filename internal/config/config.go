package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	GrpcAddress string `env:"GRPC_SERVER" envDefault:"0.0.0.0:8080"`
	HttpAddress string `env:"HTTP_SERVER" envDefault:"0.0.0.0:8081"`
	LogLevel    string `env:"LOG_LEVEL" envDefault:"DEBUG"`
	PostgresConfig
}

func NewConfig() (*Config, error) {
	config := &Config{}

	err := env.Parse(config)
	if err != nil {
		err = fmt.Errorf("config.NewConfig: %w", err)
	}
	return config, err
}

type PostgresConfig struct {
	Conn            string `env:"POSTGRES_CONN"`
	AutoMigrateUp   string `env:"AUTO_MIGRATE_UP" envDefault:"true"`
	AutoMigrateDown string `env:"AUTO_MIGRATE_DOWN" envDefault:"false"`
	MigrationsURL   string `env:"MIGRATIONS_URL" envDefault:"file:///goserver/internal/repository/db/migrations/"`
}

func NewPostgresConfig() (*PostgresConfig, error) {
	config := &PostgresConfig{}

	err := env.Parse(config)
	if err != nil {
		err = fmt.Errorf("config.NewPostgresConfig: %w", err)
	}
	return config, err
}
