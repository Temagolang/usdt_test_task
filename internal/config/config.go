package config

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Config holds all application configuration, organized by concern.
// Load priority: defaults → env vars → CLI flags (flags win).
type Config struct {
	App      App
	GRPC     GRPC
	HTTP     HTTP
	Postgres Postgres
	Grinex   Grinex
}

// App holds application-level settings.
type App struct {
	LogLevel string
}

// GRPC holds gRPC server settings.
type GRPC struct {
	Port int
}

// HTTP holds HTTP server settings (healthz endpoint).
type HTTP struct {
	Port int
}

// Postgres holds PostgreSQL connection settings.
type Postgres struct {
	DSN string
}

// Grinex holds Grinex exchange API settings.
type Grinex struct {
	URL string
}

// Load builds config by reading env vars first, then applying flag overrides.
func Load(cmd *cobra.Command) (*Config, error) {
	cfg := defaults()
	loadFromEnv(&cfg)
	applyFlags(cmd, &cfg)

	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func defaults() Config {
	return Config{
		App:  App{LogLevel: "info"},
		GRPC: GRPC{Port: 50051},
		HTTP: HTTP{Port: 8080},
		Grinex: Grinex{
			URL: "https://grinex.io",
		},
	}
}

func (c *Config) validate() error {
	if c.GRPC.Port <= 0 {
		return fmt.Errorf("invalid grpc port: %d", c.GRPC.Port)
	}
	if c.HTTP.Port <= 0 {
		return fmt.Errorf("invalid http port: %d", c.HTTP.Port)
	}
	return nil
}

// GRPCAddr returns the gRPC listen address.
func (c *Config) GRPCAddr() string {
	return fmt.Sprintf(":%d", c.GRPC.Port)
}

// HTTPAddr returns the HTTP listen address.
func (c *Config) HTTPAddr() string {
	return fmt.Sprintf(":%d", c.HTTP.Port)
}
