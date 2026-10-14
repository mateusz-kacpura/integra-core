package telemetry

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// Tracer returns an OTel tracer from the global provider.
func Tracer(name string) trace.Tracer {
	return otel.Tracer(name)
}
