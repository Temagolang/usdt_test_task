package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	transportgrpc "github.com/example/grinex-rates-service/internal/transport/grpc"
)

const shutdownTimeout = 15 * time.Second

// serve starts gRPC and HTTP servers, then blocks until a termination signal
// is received or a server fails, and performs graceful shutdown.
func (a *Application) serve(ctx context.Context, grpcSrv *transportgrpc.Server, httpSrv *http.Server) error {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// serveErr receives the first fatal error from any server goroutine.
	serveErr := make(chan error, 1)

	// Start gRPC server.
	grpcAddr := a.cfg.GRPCAddr()
	grpcLis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return fmt.Errorf("listen grpc: %w", err)
	}
	a.logger.Info("gRPC server starting", zap.String("addr", grpcAddr))
	go func() {
		if err := grpcSrv.Serve(grpcLis); err != nil {
			serveErr <- fmt.Errorf("grpc serve: %w", err)
		}
	}()

	// Bind HTTP listener synchronously to catch port conflicts early.
	httpLis, err := net.Listen("tcp", httpSrv.Addr)
	if err != nil {
		grpcSrv.GracefulStop()
		return fmt.Errorf("listen http: %w", err)
	}
	a.logger.Info("HTTP server starting", zap.String("addr", httpSrv.Addr))
	go func() {
		if err := httpSrv.Serve(httpLis); err != nil && err != http.ErrServerClosed {
			serveErr <- fmt.Errorf("http serve: %w", err)
		}
	}()

	// Wait for shutdown signal or server failure.
	var serverFailed error
	select {
	case <-ctx.Done():
		a.logger.Info("shutdown signal received")
	case serverFailed = <-serveErr:
		a.logger.Error("server failed, initiating shutdown", zap.Error(serverFailed))
		stop()
	}

	// Graceful shutdown.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	a.logger.Info("gRPC server stopping")
	grpcSrv.GracefulStop()

	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("http shutdown: %w", err)
	}

	a.logger.Info("shutdown complete")
	return serverFailed
}
