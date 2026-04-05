package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/example/grinex-rates-service/internal/app"
	"github.com/example/grinex-rates-service/internal/config"
)

var grpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Start gRPC server",
	RunE:  runGRPC,
}

func runGRPC(cmd *cobra.Command, _ []string) error {
	cfg, err := config.Load(cmd)
	if err != nil {
		return err
	}

	logger, err := buildLogger(cfg.App.LogLevel)
	if err != nil {
		return err
	}
	defer func() { _ = logger.Sync() }()

	application := app.New(cfg, logger)
	return application.Run(cmd.Context())
}

func buildLogger(level string) (*zap.Logger, error) {
	lvl, err := zapcore.ParseLevel(level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level %q: %w", level, err)
	}

	zapCfg := zap.NewProductionConfig()
	zapCfg.Level = zap.NewAtomicLevelAt(lvl)

	return zapCfg.Build()
}