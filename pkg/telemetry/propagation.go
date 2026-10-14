package telemetry

import (
	"context"
	"crypto/rand"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	json "github.com/mateusz-kacpura/integra-core/pkg/json"
)

// HeaderCorrelationID is the HTTP header used to carry the request correlation ID.
const HeaderCorrelationID = "X-Correlation-ID"

// HeaderRequestID is the HTTP header used to carry the unique request ID.
const HeaderRequestID = "X-Request-ID"

// HeaderTraceParent is the W3C trace context header for distributed trace parent information.
const HeaderTraceParent = "traceparent"

// HeaderTraceState is the W3C trace context header for vendor-specific trace state.
const HeaderTraceState = "tracestate"

type correlationKey struct{}

var wirePropagation = propagation.TraceContext{}

func newCorrelationID() string {
	return uuid.NewString()
}

func extractCorrelationID(getter func(string) string) string {
	for _, key := range []string{HeaderCorrelationID, HeaderRequestID, "correlation_id", "Correlation-Id"} {
		if val := getter(key); val != "" {
			return val
		}
	}
	return ""
}

// CorrelationID returns the correlation ID attached to ctx, or an empty string if none is set.
func CorrelationID(ctx context.Context) string {
	id, _ := ctx.Value(correlationKey{}).(string)
	return id
}

// WithCorrelationID attaches a correlation ID to the context.
// If id is empty or invalid, a new random UUID v4 is generated.
// Returns the resulting context and the validated correlation ID.
func WithCorrelationID(ctx context.Context, id string) (context.Context, string) {
	if !validID(id) {
		id = newCorrelationID()
	}
	ctx = context.WithValue(ctx, correlationKey{}, id)
	logger := log.Ctx(ctx).With().Str("correlation_id", id).Logger()
	return logger.WithContext(ctx), id
}
func validID(id string) bool {
	if len(id) == 0 || len(id) > 128 {
		return false
	}
	for _, c := range id {
		if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '-' || c == '_' || c == '.') {
			return false
		}
	}
	return true
}

// Restore extracts distributed tracing and correlation ID headers into the context, initializing span context if not present.
func Restore(ctx context.Context, headers map[string]string) context.Context {
	ctx = wirePropagation.Extract(ctx, propagation.MapCarrier(headers))
	id := headers["correlation_id"]
	if !validID(id) {
		id = newCorrelationID()
	}
	ctx = context.WithValue(ctx, correlationKey{}, id)
	if !trace.SpanContextFromContext(ctx).IsValid() {
		var tid trace.TraceID
		var sid trace.SpanID
		rand.Read(tid[:])
		rand.Read(sid[:])
		ctx = trace.ContextWithSpanContext(ctx, trace.NewSpanContext(trace.SpanContextConfig{TraceID: tid, SpanID: sid}))
	}
	logger := log.Logger.With().Str("correlation_id", id).Str("trace_id", trace.SpanContextFromContext(ctx).TraceID().String()).Logger()
	return logger.WithContext(ctx)
}

// Capture extracts correlation ID and OpenTelemetry trace metadata from the context into a map of headers.
func Capture(ctx context.Context) map[string]string {
	headers := map[string]string{}
	if id := CorrelationID(ctx); id != "" {
		headers["correlation_id"] = id
	}
	wirePropagation.Inject(ctx, propagation.MapCarrier(headers))
	return headers
}

// Metadata persists only the defined correlation/trace contract, never arbitrary request headers.
type metadataCarrier struct {
	Correlation string `json:"correlation_id,omitempty"`
	Parent      string `json:"traceparent,omitempty"`
	State       string `json:"tracestate,omitempty"`
}

func (c *metadataCarrier) Get(key string) string {
	switch key {
	case "correlation_id":
		return c.Correlation
	case "traceparent":
		return c.Parent
	case "tracestate":
		return c.State
	}
	return ""
}
func (c *metadataCarrier) Set(key, value string) {
	switch key {
	case "correlation_id":
		c.Correlation = value
	case "traceparent":
		c.Parent = value
	case "tracestate":
		c.State = value
	}
}
func (*metadataCarrier) Keys() []string {
	return []string{"correlation_id", "traceparent", "tracestate"}
}

// Metadata serializes the correlation ID and trace context from ctx into a JSON byte slice for transport or storage.
func Metadata(ctx context.Context) []byte {
	carrier := metadataCarrier{Correlation: CorrelationID(ctx)}
	wirePropagation.Inject(ctx, &carrier)
	data, _ := json.Marshal(&carrier)
	return data
}

// OutboxHeaders serializes correlation ID and trace metadata as JSON bytes for transactional outbox tables.
// If the context does not yet contain a correlation ID, a new UUID v4 is generated.
func OutboxHeaders(ctx context.Context) []byte {
	if CorrelationID(ctx) == "" {
		ctx, _ = WithCorrelationID(ctx, "")
	}
	return Metadata(ctx)
}

// DecodeMetadata unmarshals a JSON byte slice into a map of metadata headers.
func DecodeMetadata(data []byte) (map[string]string, error) {
	var headers map[string]string
	if len(data) == 0 {
		return map[string]string{}, nil
	}
	err := json.Unmarshal(data, &headers)
	return headers, err
}

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusResponseWriter) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// HTTP wraps an http.Handler with distributed tracing, correlation extraction, response header injection, and access logging.
func HTTP(next http.Handler) http.Handler {
	tracer := Tracer("integra/http")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		id := extractCorrelationID(r.Header.Get)
		ctx := Restore(r.Context(), map[string]string{
			"correlation_id":  id,
			HeaderTraceParent: r.Header.Get(HeaderTraceParent),
			HeaderTraceState:  r.Header.Get(HeaderTraceState),
		})
		corrID := CorrelationID(ctx)
		w.Header().Set(HeaderCorrelationID, corrID)
		w.Header().Set(HeaderRequestID, corrID)

		ctx, span := tracer.Start(ctx, r.Method+" "+r.URL.Path,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				attribute.String("correlation_id", corrID),
				attribute.String("http.method", r.Method),
				attribute.String("http.path", r.URL.Path),
			),
		)
		defer span.End()

		sw := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(sw, r.WithContext(ctx))

		span.SetAttributes(attribute.Int("http.status_code", sw.statusCode))
		if sw.statusCode >= 400 {
			span.SetStatus(codes.Error, http.StatusText(sw.statusCode))
		} else {
			span.SetStatus(codes.Ok, "")
		}

		log.Ctx(ctx).Info().
			Int("status", sw.statusCode).
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Dur("duration", time.Since(start)).
			Str("correlation_id", corrID).
			Str("trace_id", span.SpanContext().TraceID().String()).
			Str("span_id", span.SpanContext().SpanID().String()).
			Msg("http access log")
	})
}

// RoundTripper is an http.RoundTripper that automatically propagates the
// correlation ID and W3C trace context from context.Context to outgoing HTTP request headers.
type RoundTripper struct {
	Next http.RoundTripper
}

// NewRoundTripper wraps the provided RoundTripper (or http.DefaultTransport if nil)
// to inject X-Correlation-ID, X-Request-ID, and W3C traceparent into downstream HTTP requests.
func NewRoundTripper(next http.RoundTripper) http.RoundTripper {
	if next == nil {
		next = http.DefaultTransport
	}
	return &RoundTripper{Next: next}
}

// RoundTrip executes a single HTTP transaction, injecting correlation and trace headers from ctx.
func (rt *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	id := CorrelationID(ctx)
	if id == "" {
		ctx, id = WithCorrelationID(ctx, "")
	}
	ctx, span := Tracer("integra/http-client").Start(ctx, req.Method+" "+req.URL.Path,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("correlation_id", id),
			attribute.String("http.method", req.Method),
			attribute.String("http.url", req.URL.String()),
		),
	)
	defer span.End()

	out := req.Clone(ctx)
	if out.Header.Get(HeaderCorrelationID) == "" {
		out.Header.Set(HeaderCorrelationID, id)
	}
	if out.Header.Get(HeaderRequestID) == "" {
		out.Header.Set(HeaderRequestID, id)
	}
	wirePropagation.Inject(ctx, propagation.HeaderCarrier(out.Header))
	resp, err := rt.Next.RoundTrip(out)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	} else {
		span.SetAttributes(attribute.Int("http.status_code", resp.StatusCode))
		if resp.StatusCode >= 400 {
			span.SetStatus(codes.Error, http.StatusText(resp.StatusCode))
		} else {
			span.SetStatus(codes.Ok, "")
		}
	}
	return resp, err
}

// NewHTTPClient returns an http.Client equipped with downstream correlation propagation.
func NewHTTPClient(base *http.Client) *http.Client {
	if base == nil {
		return &http.Client{Transport: NewRoundTripper(nil)}
	}
	cloned := *base
	cloned.Transport = NewRoundTripper(base.Transport)
	return &cloned
}

// InjectHTTP writes the correlation ID and trace headers from ctx into the provided HTTP header.
func InjectHTTP(ctx context.Context, header http.Header) {
	id := CorrelationID(ctx)
	if id == "" {
		return
	}
	if header.Get(HeaderCorrelationID) == "" {
		header.Set(HeaderCorrelationID, id)
	}
	if header.Get(HeaderRequestID) == "" {
		header.Set(HeaderRequestID, id)
	}
	wirePropagation.Inject(ctx, propagation.HeaderCarrier(header))
}
