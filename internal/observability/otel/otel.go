// Package otel provides OpenTelemetry initialization and shutdown.
// This is a skeleton foundation — full tracing of Grinex/Postgres
// will be added in a later step.
package otel

import (
	"context"

	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// ShutdownFunc gracefully shuts down the OTel provider.
type ShutdownFunc func(ctx context.Context) error

// Init sets up a basic TracerProvider.
// Returns a shutdown function that must be called on application exit.
func Init(_ context.Context) (ShutdownFunc, error) {
	tp := sdktrace.NewTracerProvider()
	otel.SetTracerProvider(tp)

	shutdown := func(ctx context.Context) error {
		return tp.Shutdown(ctx)
	}
	return shutdown, nil
}
