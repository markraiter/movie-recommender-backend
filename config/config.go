package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v8"
)

type Server struct {
	AppAddress      string        `env:"APP_PORT" envDefault:"8000"`
	AppReadTimeout  time.Duration `env:"APP_READ_TIMEOUT" envDefault:"60s"`
	AppWriteTimeout time.Duration `env:"APP_WRITE_TIMEOUT" envDefault:"60s"`
	AppIdleTimeout  time.Duration `env:"APP_IDLE_TIMEOUT" envDefault:"60s"`
}

type Mongo struct {
	ConnectionString string `env:"MONGO_CONNECTION_STRING" envDefault:"mongodb://localhost:27017"`
	NameDB           string `env:"MONGO_DB" envDefault:"movie-recommender"`
	Username         string `env:"MONGO_USERNAME" envDefault:"root"`
	Password         string `env:"MONGO_PASSWORD" envDefault:"example"`
}

type Config struct {
	Server   Server
	DB       Mongo
	Auth     AuthConfig
	LogLevel string `env:"LOG_LEVEL" envDefault:"INFO"`
}

type AuthConfig struct {
	Salt            string        `env:"APP_SALT,notEmpty"`
	SigningKey      string        `env:"SIGNING_KEY,notEmpty"`
	AccessTokenTTL  time.Duration `env:"ACCESS_TOKEN_TTL" envDefault:"15m"`
	RefreshTokenTTL time.Duration `env:"REFRESH_TOKEN_TTL" envDefault:"24h"`
}

func InitConfig() (Config, error) {
	cfg := Config{}

	if err := env.Parse(&cfg); err != nil {
		return Config{}, fmt.Errorf("error while parsing .env: %w", err)
	}

	cfg.Server.AppAddress = ":" + cfg.Server.AppAddress

	return cfg, nil
}
