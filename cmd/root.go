package cmd

import (
	"github.com/spf13/cobra"

	"github.com/example/grinex-rates-service/internal/config"
)

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "USDT rates service with Grinex integration",
	// Bare ./app starts the gRPC server (acceptance flow compatibility).
	RunE: runGRPC,
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	config.RegisterFlags(rootCmd)

	rootCmd.AddCommand(grpcCmd)
	rootCmd.AddCommand(migrateCmd)
}
