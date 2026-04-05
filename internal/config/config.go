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

type App struct {
	LogLevel string
}

type GRPC struct {
	Port int
}

type HTTP struct {
	Port int
}

type Postgres struct {
	DSN string
}

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