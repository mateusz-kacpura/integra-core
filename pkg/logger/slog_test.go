package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mateusz-kacpura/integra-core/pkg/telemetry"
)

func TestMultiHandlerBroadcast(t *testing.T) {
	buf1 := &bytes.Buffer{}
	buf2 := &bytes.Buffer{}

	h1 := slog.NewJSONHandler(buf1, &slog.HandlerOptions{Level: slog.LevelInfo})
	h2 := slog.NewJSONHandler(buf2, &slog.HandlerOptions{Level: slog.LevelInfo})

	multi := NewMultiHandler(h1, h2)
	l := slog.New(multi)

	l.Info("broadcast test message")

	var r1, r2 map[string]any
	require.NoError(t, json.Unmarshal(buf1.Bytes(), &r1))
	require.NoError(t, json.Unmarshal(buf2.Bytes(), &r2))

	assert.Equal(t, "broadcast test message", r1["msg"])
	assert.Equal(t, "broadcast test message", r2["msg"])
}

func TestOperationalAndAuditHandlerStreamSeparation(t *testing.T) {
	appBuf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}

	handler := NewOperationalAndAuditHandler(appBuf, errBuf, slog.LevelInfo)
	l := slog.New(handler)

	l.Info("normal operational log")
	l.Error("critical system failure")

	// Info should go to appBuf only
	assert.Contains(t, appBuf.String(), "normal operational log")
	assert.NotContains(t, appBuf.String(), "critical system failure")

	// Error should go to errBuf only
	assert.Contains(t, errBuf.String(), "critical system failure")
	assert.NotContains(t, errBuf.String(), "normal operational log")
}

func TestTraceAndSanitizingHandlerPIIRedaction(t *testing.T) {
	buf := &bytes.Buffer{}
	baseHandler := slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelInfo})
	handler := NewTraceAndSanitizingHandler(baseHandler)
	l := slog.New(handler)

	ctx, corrID := telemetry.WithCorrelationID(context.Background(), "corr-trace-test")

	// Message with credit card and email PII
	l.InfoContext(ctx, "Customer 4111 1111 1111 1234 registered with email user@domain.com")

	var r map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &r))

	assert.Equal(t, corrID, r["correlation_id"])
	msg, ok := r["msg"].(string)
	require.True(t, ok)
	assert.NotContains(t, msg, "4111 1111 1111 1234")
	assert.NotContains(t, msg, "user@domain.com")
	assert.Contains(t, msg, "4111 **** **** 1234")
	assert.Contains(t, msg, "u***@domain.com")
}

func TestNewSlogLogger(t *testing.T) {
	l := NewSlogLogger(EntrypointConfig{
		Service:     "iam-service",
		ProcessType: "server",
		Environment: "production",
		Version:     "v1.0.0",
		LogLevel:    "info",
	})
	assert.NotNil(t, l)
}
