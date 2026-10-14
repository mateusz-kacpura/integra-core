// Package telemetry provides OpenTelemetry initialization, distributed tracing propagation,
// Prometheus metric collectors, and zerolog integration for microservices.
//
// Usage:
//
//	shutdown, err := telemetry.InitConditional(ctx, "my-service", true, "otel-collector:4318")
//	defer shutdown(ctx)
//
//	logger := zerolog.New(os.Stdout).Hook(telemetry.NewTraceHook())
package telemetry
