package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"

	"go.opentelemetry.io/otel/trace"

	"github.com/mateusz-kacpura/integra-core/pkg/telemetry"
)

// MultiHandler dispatches log records to multiple slog.Handlers.
type MultiHandler struct {
	handlers []slog.Handler
}

// NewMultiHandler creates a new MultiHandler combining multiple slog handlers.
func NewMultiHandler(handlers ...slog.Handler) slog.Handler {
	return &MultiHandler{handlers: handlers}
}

// Enabled returns true if any underlying handler is enabled for the level.
func (m *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

// Handle forwards the record to all enabled handlers.
func (m *MultiHandler) Handle(ctx context.Context, r slog.Record) error {
	var firstErr error
	for _, h := range m.handlers {
		if h.Enabled(ctx, r.Level) {
			if err := h.Handle(ctx, r.Clone()); err != nil && firstErr == nil {
				firstErr = err
			}
		}
	}
	return firstErr
}

// WithAttrs returns a new MultiHandler with attributes added to all child handlers.
func (m *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithAttrs(attrs)
	}
	return &MultiHandler{handlers: handlers}
}

// WithGroup returns a new MultiHandler with a group name added to all child handlers.
func (m *MultiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithGroup(name)
	}
	return &MultiHandler{handlers: handlers}
}

// LevelFilterHandler filters records based on minimum and optional maximum slog.Level.
type LevelFilterHandler struct {
	handler  slog.Handler
	minLevel slog.Level
	maxLevel *slog.Level
}

// NewLevelFilterHandler creates a handler filtering by minimum level.
func NewLevelFilterHandler(h slog.Handler, minLevel slog.Level) slog.Handler {
	return &LevelFilterHandler{handler: h, minLevel: minLevel}
}

// NewLevelRangeHandler creates a handler filtering within [minLevel, maxLevel].
func NewLevelRangeHandler(h slog.Handler, minLevel slog.Level, maxLevel slog.Level) slog.Handler {
	return &LevelFilterHandler{handler: h, minLevel: minLevel, maxLevel: &maxLevel}
}

// Enabled returns true if the level is within the configured range.
func (f *LevelFilterHandler) Enabled(ctx context.Context, level slog.Level) bool {
	if level < f.minLevel {
		return false
	}
	if f.maxLevel != nil && level > *f.maxLevel {
		return false
	}
	return f.handler.Enabled(ctx, level)
}

// Handle executes the wrapped handler if the level is accepted.
func (f *LevelFilterHandler) Handle(ctx context.Context, r slog.Record) error {
	if !f.Enabled(ctx, r.Level) {
		return nil
	}
	return f.handler.Handle(ctx, r)
}

// WithAttrs returns a new LevelFilterHandler with attributes attached.
func (f *LevelFilterHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &LevelFilterHandler{
		handler:  f.handler.WithAttrs(attrs),
		minLevel: f.minLevel,
		maxLevel: f.maxLevel,
	}
}

// WithGroup returns a new LevelFilterHandler with a group attached.
func (f *LevelFilterHandler) WithGroup(name string) slog.Handler {
	return &LevelFilterHandler{
		handler:  f.handler.WithGroup(name),
		minLevel: f.minLevel,
		maxLevel: f.maxLevel,
	}
}

// TraceAndSanitizingHandler wraps an slog.Handler to inject OpenTelemetry trace/span IDs,
// correlation IDs, and sanitize PII in message strings and string attributes.
type TraceAndSanitizingHandler struct {
	next slog.Handler
}

// NewTraceAndSanitizingHandler wraps a handler with trace correlation and PII sanitization.
func NewTraceAndSanitizingHandler(h slog.Handler) slog.Handler {
	return &TraceAndSanitizingHandler{next: h}
}

// Enabled delegates to the wrapped handler.
func (h *TraceAndSanitizingHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}

// Handle injects correlation and trace attributes and sanitizes message text.
func (h *TraceAndSanitizingHandler) Handle(ctx context.Context, r slog.Record) error {
	r.Message = SanitizeString(r.Message)

	var extraAttrs []slog.Attr
	if corrID := telemetry.CorrelationID(ctx); corrID != "" {
		extraAttrs = append(extraAttrs, slog.String("correlation_id", corrID))
	}
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		extraAttrs = append(extraAttrs,
			slog.String("trace_id", span.SpanContext().TraceID().String()),
			slog.String("span_id", span.SpanContext().SpanID().String()),
		)
	}

	if len(extraAttrs) > 0 {
		r.AddAttrs(extraAttrs...)
	}

	return h.next.Handle(ctx, r)
}

// WithAttrs sanitizes string attributes before attaching them.
func (h *TraceAndSanitizingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	sanitized := make([]slog.Attr, len(attrs))
	for i, a := range attrs {
		if a.Value.Kind() == slog.KindString {
			sanitized[i] = slog.String(a.Key, SanitizeString(a.Value.String()))
		} else {
			sanitized[i] = a
		}
	}
	return &TraceAndSanitizingHandler{next: h.next.WithAttrs(sanitized)}
}

// WithGroup delegates to the wrapped handler.
func (h *TraceAndSanitizingHandler) WithGroup(name string) slog.Handler {
	return &TraceAndSanitizingHandler{next: h.next.WithGroup(name)}
}

// NewOperationalAndAuditHandler creates a composite multi-handler separating
// operational logs (stdout) from critical errors (stderr) and security events.
func NewOperationalAndAuditHandler(appWriter, errWriter io.Writer, minLevel slog.Level) slog.Handler {
	if appWriter == nil {
		appWriter = os.Stdout
	}
	if errWriter == nil {
		errWriter = os.Stderr
	}

	sanitizedApp := NewSanitizingWriter(appWriter)
	sanitizedErr := NewSanitizingWriter(errWriter)

	warnLevel := slog.LevelWarn
	appHandler := NewLevelRangeHandler(
		slog.NewJSONHandler(sanitizedApp, &slog.HandlerOptions{Level: minLevel}),
		minLevel,
		warnLevel,
	)

	errHandler := NewLevelFilterHandler(
		slog.NewJSONHandler(sanitizedErr, &slog.HandlerOptions{Level: slog.LevelError}),
		slog.LevelError,
	)

	multi := NewMultiHandler(appHandler, errHandler)
	return NewTraceAndSanitizingHandler(multi)
}

// NewSlogLogger creates an *slog.Logger configured with operational context,
// multi-stream separation (stdout for app logs, stderr for errors), trace correlation,
// and PII sanitization.
func NewSlogLogger(cfg EntrypointConfig) *slog.Logger {
	level := slog.LevelInfo
	switch strings.ToLower(cfg.LogLevel) {
	case "debug":
		level = slog.LevelDebug
	case "warn", "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	handler := NewOperationalAndAuditHandler(os.Stdout, os.Stderr, level)

	var metaAttrs []slog.Attr
	if cfg.Service != "" {
		metaAttrs = append(metaAttrs, slog.String("service", cfg.Service), slog.String("app.name", cfg.Service))
	}
	proc := cfg.ProcessType
	if proc == "" {
		proc = cfg.Component
	}
	if proc != "" {
		metaAttrs = append(metaAttrs, slog.String("process_type", proc), slog.String("component", proc))
	}
	if cfg.Environment != "" {
		metaAttrs = append(metaAttrs, slog.String("environment", cfg.Environment))
	}
	if cfg.Version != "" {
		metaAttrs = append(metaAttrs, slog.String("version", cfg.Version))
	}

	if len(metaAttrs) > 0 {
		handler = handler.WithAttrs(metaAttrs)
	}

	return slog.New(handler)
}
