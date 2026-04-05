package app

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"

	"github.com/example/grinex-rates-service/internal/client/grinex"
	"github.com/example/grinex-rates-service/internal/config"
	appmetrics "github.com/example/grinex-rates-service/internal/observability/metrics"
	otelinit "github.com/example/grinex-rates-service/internal/observability/otel"
	postgresrepo "github.com/example/grinex-rates-service/internal/repo/rates/postgres"
	"github.com/example/grinex-rates-service/internal/service/rates"
	transportgrpc "github.com/example/grinex-rates-service/internal/transport/grpc"
)

const (
	otelShutdownTimeout = 5 * time.Second

	// Grinex trading symbol for USDT/RUB depth.
	// Hardcoded per spec: the service is specifically for this pair.
	grinexSymbol = "usdta7a5"
)

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

	// Prometheus metrics.
	promRegistry := prometheus.NewRegistry()
	metrics := appmetrics.New(promRegistry)

	// Connect to PostgreSQL.
	pool, err := pgxpool.New(ctx, a.cfg.Postgres.DSN)
	if err != nil {
		return fmt.Errorf("connect to postgres: %w", err)
	}
	defer pool.Close()

	// Wire dependencies.
	grinexClient := grinex.New(a.cfg.Grinex.URL, grinexSymbol, metrics)
	repo := postgresrepo.New(pool)
	svc := rates.NewService(grinexClient, repo)

	// Create servers.
	grpcSrv := transportgrpc.NewServer(svc, metrics, a.logger)
	httpSrv := newHTTPServer(a.cfg.HTTPAddr(), promRegistry)

	// Start servers and wait for graceful shutdown.
	return a.serve(ctx, grpcSrv, httpSrv)
}
