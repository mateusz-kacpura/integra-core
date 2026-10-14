package telemetry

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.opentelemetry.io/otel/trace"

	json "github.com/mateusz-kacpura/integra-core/pkg/json"
)

func TestRequestMetadataSurvivesPersistenceAndIndependentConsumer(t *testing.T) {
	parent := "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01"
	request := httptest.NewRequest("POST", "/api", nil)
	request.Header.Set("X-Correlation-ID", "registration-123")
	request.Header.Set("traceparent", parent)
	request.Header.Set("tracestate", "vendor=opaque")
	var stored []byte
	response := httptest.NewRecorder()
	HTTP(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) { stored = Metadata(r.Context()) })).ServeHTTP(response, request)
	if response.Header().Get("X-Correlation-ID") != "registration-123" {
		t.Fatal(response.Header())
	}
	headers, err := DecodeMetadata(stored)
	if err != nil {
		t.Fatal(err)
	}
	consumer := Restore(context.Background(), headers)
	if CorrelationID(consumer) != "registration-123" {
		t.Fatal("lost correlation")
	}
	if trace.SpanContextFromContext(consumer).TraceID().String() != "4bf92f3577b34da6a3ce929d0e0e4736" {
		t.Fatal("lost trace")
	}
	if Capture(consumer)["tracestate"] != "vendor=opaque" {
		t.Fatal("lost tracestate")
	}
	retry := Restore(context.Background(), headers)
	if string(Metadata(retry)) != string(Metadata(consumer)) {
		t.Fatal("retry changed durable metadata")
	}
}
func TestRequestMetadataRejectsUnboundedAndUnsafeIDs(t *testing.T) {
	for _, value := range []string{"", "bad\nheader", string(make([]byte, 129))} {
		ctx := Restore(context.Background(), map[string]string{"correlation_id": value, "traceparent": "invalid"})
		if !validID(CorrelationID(ctx)) || !trace.SpanContextFromContext(ctx).IsValid() {
			t.Fatal("invalid generated context")
		}
	}
}

func TestTypedMetadataMatchesPropagationContract(t *testing.T) {
	for _, ctx := range []context.Context{context.Background(), Restore(context.Background(), map[string]string{"correlation_id": "test", "traceparent": "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-03", "tracestate": "vendor=value"})} {
		expected, err := json.Marshal(Capture(ctx))
		if err != nil {
			t.Fatal(err)
		}
		if got := string(Metadata(ctx)); got != string(expected) {
			t.Fatalf("propagation changed: %s vs %s", got, expected)
		}
	}
}

func TestOutboxHeadersEnsuresCorrelationID(t *testing.T) {
	emptyCtx := context.Background()
	outboxBytes := OutboxHeaders(emptyCtx)
	decoded, err := DecodeMetadata(outboxBytes)
	if err != nil {
		t.Fatalf("failed to decode outbox headers: %v", err)
	}
	corrID := decoded["correlation_id"]
	if !validID(corrID) || corrID == "" {
		t.Fatalf("expected valid generated correlation ID, got %q", corrID)
	}

	existingID := "saga-transaction-456"
	ctx, _ := WithCorrelationID(context.Background(), existingID)
	outboxBytesExisting := OutboxHeaders(ctx)
	decodedExisting, err := DecodeMetadata(outboxBytesExisting)
	if err != nil {
		t.Fatalf("failed to decode outbox headers: %v", err)
	}
	if decodedExisting["correlation_id"] != existingID {
		t.Fatalf("expected %q, got %q", existingID, decodedExisting["correlation_id"])
	}
}
