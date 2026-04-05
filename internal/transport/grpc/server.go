package grpc

import (
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/example/grinex-rates-service/internal/service/rates"
)

// Server wraps a gRPC server with registered services.
type Server struct {
	srv    *grpc.Server
	health *health.Server
}

// NewServer creates a gRPC server with health and rates services registered.
func NewServer(svc *rates.Service, logger *zap.Logger) *Server {
	_ = logger // will be used for interceptors in T16

	srv := grpc.NewServer()

	// Register standard gRPC health service.
	healthSrv := health.NewServer()
	healthpb.RegisterHealthServer(srv, healthSrv)

	// Register rates service handler when available.
	if svc != nil {
		NewRatesHandler(srv, svc)
	}

	// Enable server reflection for development.
	reflection.Register(srv)

	return &Server{
		srv:    srv,
		health: healthSrv,
	}
}

// Serve starts serving on the given listener.
func (s *Server) Serve(lis net.Listener) error {
	return s.srv.Serve(lis)
}

// GracefulStop gracefully stops the gRPC server.
func (s *Server) GracefulStop() {
	s.health.Shutdown()
	s.srv.GracefulStop()
}
