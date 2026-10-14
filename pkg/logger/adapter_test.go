package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mateusz-kacpura/integra-core/pkg/telemetry"
)

func TestPortAdapterSanitizesAndPropagatesTrace(t *testing.T) {
	buf := &bytes.Buffer{}
	zl := zerolog.New(buf).With().Timestamp().Logger()
	adapter := NewPortAdapter(zl)

	ctx, corrID := telemetry.WithCorrelationID(context.Background(), "corr-adapter-test")

	adapter.LogCtx(ctx, "warn", "User 4111 1111 1111 9999 failed login", map[string]any{
		"secret_email":  "alice@company.com",
		"attempt_count": 3,
	})

	var r map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &r))

	assert.Equal(t, "warn", r["level"])
	assert.Equal(t, corrID, r["correlation_id"])

	msg, ok := r["message"].(string)
	require.True(t, ok)
	assert.NotContains(t, msg, "4111 1111 1111 9999")
	assert.Contains(t, msg, "4111 **** **** 9999")

	secretEmail, ok := r["secret_email"].(string)
	require.True(t, ok)
	assert.NotContains(t, secretEmail, "alice@company.com")
	assert.Contains(t, secretEmail, "a***@company.com")
	assert.Equal(t, float64(3), r["attempt_count"])
}

func TestPortAdapterLogWithoutContext(t *testing.T) {
	buf := &bytes.Buffer{}
	zl := zerolog.New(buf).With().Timestamp().Logger()
	adapter := NewPortAdapter(zl)

	adapter.Log("error", "database connection failed", map[string]any{"timeout": "5s"})

	var r map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &r))

	assert.Equal(t, "error", r["level"])
	assert.Equal(t, "database connection failed", r["message"])
}
