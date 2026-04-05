package config

import "github.com/spf13/cobra"

// RegisterFlags adds persistent CLI flags so subcommands inherit them.
func RegisterFlags(cmd *cobra.Command) {
	f := cmd.PersistentFlags()
	f.Int("grpc-port", 50051, "gRPC server port")
	f.Int("http-port", 8080, "HTTP server port")
	f.String("database-dsn", "", "PostgreSQL connection string")
	f.String("log-level", "info", "Log level (debug, info, warn, error)")
	f.String("grinex-url", "", "Grinex API base URL")
}

func applyFlags(cmd *cobra.Command, cfg *Config) {
	if cmd.Flags().Changed("grpc-port") {
		cfg.GRPC.Port, _ = cmd.Flags().GetInt("grpc-port")
	}
	if cmd.Flags().Changed("http-port") {
		cfg.HTTP.Port, _ = cmd.Flags().GetInt("http-port")
	}
	if cmd.Flags().Changed("database-dsn") {
		cfg.Postgres.DSN, _ = cmd.Flags().GetString("database-dsn")
	}
	if cmd.Flags().Changed("log-level") {
		cfg.App.LogLevel, _ = cmd.Flags().GetString("log-level")
	}
	if cmd.Flags().Changed("grinex-url") {
		cfg.Grinex.URL, _ = cmd.Flags().GetString("grinex-url")
	}
}
