package config

import (
	"fmt"
	"regexp"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	GrpcAddress string `env:"GRPC_SERVER" envDefault:"0.0.0.0:8080"`
	GrpcPort    string
	HttpAddress string `env:"HTTP_SERVER" envDefault:"0.0.0.0:8081"`
	HttpPort    string
	LogLevel    string `env:"LOG_LEVEL" envDefault:"DEBUG"`
	PostgresConfig
}

func NewConfig() (*Config, error) {
	config := &Config{}

	err := env.Parse(config)
	if err != nil {
		err = fmt.Errorf("config.NewConfig: %w", err)
	}

	re := regexp.MustCompile(`:(\d)*`)
	config.GrpcPort = re.FindString(config.GrpcAddress)[1:]
	config.HttpPort = re.FindString(config.HttpAddress)[1:]

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
