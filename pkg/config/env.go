package config

import (
	"time"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Port     string        `env:"PORT,required"`
	DBHost   string        `env:"DB_HOST,required"`
	DBPort   string        `env:"DB_PORT,required"`
	DBUser   string        `env:"DB_USER,required"`
	DBPass   string        `env:"DB_PASS,required"`
	DBName   string        `env:"DB_NAME,required"`
	SSLMode  bool          `env:"DB_SSL_MODE,required"`
	TokenTTL time.Duration `env:"TOKEN_TTL,required"`
	CacheTTL time.Duration `env:"CACHE_TTL,required"`
}

func Parse() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
