package logger

import (
	"context"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/trace"

	"github.com/mateusz-kacpura/integra-core/pkg/telemetry"
)

// PortAdapter adapts zerolog.Logger to application port Logger interfaces,
// ensuring PII sanitization and OpenTelemetry trace propagation.
type PortAdapter struct {
	Logger zerolog.Logger
}

// NewPortAdapter creates a PortAdapter wrapping the provided zerolog.Logger.
func NewPortAdapter(zl zerolog.Logger) PortAdapter {
	return PortAdapter{Logger: zl}
}

// Log logs a message at the given level with fields without explicit context.
func (a PortAdapter) Log(level, message string, fields map[string]any) {
	a.LogCtx(context.Background(), level, message, fields)
}

// LogCtx logs a message with context, attaching trace/correlation IDs and sanitizing PII.
func (a PortAdapter) LogCtx(ctx context.Context, level, message string, fields map[string]any) {
	sanitizedMsg := SanitizeString(message)

	l := a.Logger
	if ctxLogger := log.Ctx(ctx); ctxLogger != nil && ctxLogger.GetLevel() != zerolog.Disabled {
		l = *ctxLogger
	}

	var event *zerolog.Event
	switch strings.ToLower(level) {
	case "error", "err":
		event = l.Error()
	case "warn", "warning":
		event = l.Warn()
	case "debug":
		event = l.Debug()
	default:
		event = l.Info()
	}

	if corrID := telemetry.CorrelationID(ctx); corrID != "" {
		event = event.Str("correlation_id", corrID)
	}
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		event = event.Str("trace_id", span.SpanContext().TraceID().String()).
			Str("span_id", span.SpanContext().SpanID().String())
	}

	if len(fields) > 0 {
		sanitizedFields := make(map[string]any, len(fields))
		for k, v := range fields {
			if s, ok := v.(string); ok {
				sanitizedFields[k] = SanitizeString(s)
			} else {
				sanitizedFields[k] = v
			}
		}
		event = event.Fields(sanitizedFields)
	}

	event.Msg(sanitizedMsg)
}
