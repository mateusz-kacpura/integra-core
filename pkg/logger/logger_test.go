package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mateusz-kacpura/integra-core/pkg/telemetry"
)

func TestNewWithTraceHook(t *testing.T) {
	l := NewWithTraceHook("test-service")
	// zerolog.Logger is a value type; a non-zero logger should be usable.
	// Verify it can write without panicking.
	l.Info().Msg("trace hook logger works")
}

func TestNewWithSanitizer(t *testing.T) {
	l := NewWithSanitizer("test-service")
	l.Info().Msg("sanitizer logger works")
}

func TestInitLoggerLocalDev(t *testing.T) {
	l := InitLogger("debug", true)
	l.Debug().Msg("local dev logger")
}

func TestInitLoggerProduction(t *testing.T) {
	l := InitLogger("info", false)
	l.Info().Msg("production logger")
}

func TestInitLoggerInvalidLevel(t *testing.T) {
	// Invalid level string should fall back to info level without error.
	l := InitLogger("not-a-level", false)
	l.Info().Msg("fallback logger")
}

func TestInitEntrypointWithWriter(t *testing.T) {
	buf := &bytes.Buffer{}
	ctx, l := InitEntrypointWithWriter(context.Background(), EntrypointConfig{
		Service:     "iam-service",
		Component:   "server",
		Environment: "staging",
		Version:     "v0.2.0-test.1",
	}, buf)

	assert.NotNil(t, ctx)
	corrID := telemetry.CorrelationID(ctx)
	assert.NotEmpty(t, corrID)

	l.Info().Msg("application started")

	var output map[string]any
	err := json.Unmarshal(buf.Bytes(), &output)
	assert.NoError(t, err)

	assert.Equal(t, "info", output["level"])
	assert.Equal(t, "iam-service", output["service"])
	assert.Equal(t, "iam-service", output["app.name"])
	assert.Equal(t, "server", output["component"])
	assert.Equal(t, "server", output["process_type"])
	assert.Equal(t, "staging", output["environment"])
	assert.Equal(t, "v0.2.0-test.1", output["version"])
	assert.Equal(t, corrID, output["correlation_id"])
	assert.Equal(t, "application started", output["message"])
}

func TestInitEntrypointDefaults(t *testing.T) {
	buf := &bytes.Buffer{}
	ctx, l := InitEntrypointWithWriter(context.Background(), EntrypointConfig{
		ProcessType: "worker",
	}, buf)

	assert.NotNil(t, ctx)
	l.Error().Msg("worker processing failed")

	var output map[string]any
	err := json.Unmarshal(buf.Bytes(), &output)
	assert.NoError(t, err)

	assert.Equal(t, "error", output["level"])
	assert.Equal(t, "worker", output["component"])
	assert.Equal(t, "worker", output["process_type"])
	assert.Equal(t, "development", output["environment"])
	assert.Equal(t, "v0.2.0-dev.0", output["version"])
	assert.NotEmpty(t, output["correlation_id"])
}

func TestNewBootstrapLogger(t *testing.T) {
	l := NewBootstrapLogger(EntrypointConfig{
		Service:     "ledger-service",
		ProcessType: "server",
		Environment: "production",
		LogLevel:    "debug",
	})
	assert.NotNil(t, l)
	assert.Equal(t, "debug", l.GetLevel().String())
}
