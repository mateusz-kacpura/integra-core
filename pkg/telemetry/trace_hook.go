package telemetry

import (
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

// TraceHook is a zerolog Hook that adds correlation_id, trace_id, and span_id from the context.
type TraceHook struct{}

// NewTraceHook creates a new TraceHook
func NewTraceHook() TraceHook {
	return TraceHook{}
}

// Run implements zerolog.Hook
func (h TraceHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	ctx := e.GetCtx()
	if ctx == nil {
		return
	}

	if corrID := CorrelationID(ctx); corrID != "" {
		e.Str("correlation_id", corrID)
	}

	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return
	}

	e.Str("trace_id", span.SpanContext().TraceID().String())
	e.Str("span_id", span.SpanContext().SpanID().String())
}
