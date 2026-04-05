package app

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/example/grinex-rates-service/internal/config"
	otelinit "github.com/example/grinex-rates-service/internal/observability/otel"
	transportgrpc "github.com/example/grinex-rates-service/internal/transport/grpc"
)

const otelShutdownTimeout = 5 * time.Second

// Application is the composition root that wires all dependencies
// and manages server lifecycle.
type Application struct {
	cfg    *config.Config
	logger *zap.Logger
}

// New creates a new Application.
func New(cfg *config.Config, logger *zap.Logger) *Application {
	return &Application{
		cfg:    cfg,
		logger: logger,
	}
}

// Run starts all servers and blocks until shutdown.
func (a *Application) Run(ctx context.Context) error {
	// Initialize OpenTelemetry.
	shutdownOTel, err := otelinit.Init(ctx)
	if err != nil {
		return err
	}
	defer func() {
		otelCtx, cancel := context.WithTimeout(context.Background(), otelShutdownTimeout)
		defer cancel()
		if otelErr := shutdownOTel(otelCtx); otelErr != nil {
			a.logger.Error("otel shutdown error", zap.Error(otelErr))
		}
	}()

	// Create gRPC server.
	// TODO(T15): pass real rates.Service instead of nil after wiring.
	grpcSrv := transportgrpc.NewServer(nil, a.logger)

	// Create HTTP server for /healthz.
	httpSrv := newHTTPServer(a.cfg.HTTPAddr())

	// Start servers and wait for graceful shutdown.
	return a.serve(ctx, grpcSrv, httpSrv)
}
