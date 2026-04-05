package config

import (
	"os"
	"strconv"
)

func loadFromEnv(cfg *Config) {
	if v := os.Getenv("GRPC_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			cfg.GRPC.Port = port
		}
	}
	if v := os.Getenv("HTTP_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			cfg.HTTP.Port = port
		}
	}
	if v := os.Getenv("DATABASE_DSN"); v != "" {
		cfg.Postgres.DSN = v
	}
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		cfg.App.LogLevel = v
	}
	if v := os.Getenv("GRINEX_URL"); v != "" {
		cfg.Grinex.URL = v
	}
}
