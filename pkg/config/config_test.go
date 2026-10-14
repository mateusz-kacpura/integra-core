package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadBaseConfigDefaults(t *testing.T) {
	// LoadBaseConfig uses viper defaults for SERVICE_NAME, ENV, HTTP_PORT, LOG_LEVEL.
	// Without any env overrides the defaults should apply.
	cfg, err := LoadBaseConfig()
	require.NoError(t, err)
	assert.Equal(t, "integra-service", cfg.ServiceName)
	assert.Equal(t, "development", cfg.Environment)
	assert.Equal(t, 8080, cfg.HTTPPort)
	assert.Equal(t, "info", cfg.LogLevel)
}

func TestLoadBaseConfigEnvOverrides(t *testing.T) {
	// Viper with SetEnvPrefix + AutomaticEnv only picks up env vars for keys
	// that have been explicitly bound or have defaults. LoadBaseConfig sets
	// defaults for SERVICE_NAME, ENV, HTTP_PORT, LOG_LEVEL — so those are
	// the fields we can override via env.
	t.Setenv("INTEGRA_SERVICE_NAME", "test-svc")
	t.Setenv("INTEGRA_ENV", "production")
	t.Setenv("INTEGRA_HTTP_PORT", "9090")
	t.Setenv("INTEGRA_LOG_LEVEL", "debug")

	cfg, err := LoadBaseConfig()
	require.NoError(t, err)
	assert.Equal(t, "test-svc", cfg.ServiceName)
	assert.Equal(t, "production", cfg.Environment)
	assert.Equal(t, 9090, cfg.HTTPPort)
	assert.Equal(t, "debug", cfg.LogLevel)
}

func TestValidateSuccess(t *testing.T) {
	cfg := BaseConfig{
		ServiceName: "my-service",
		DatabaseURL: "postgres://localhost/db",
		HTTPPort:    8080,
	}
	assert.NoError(t, cfg.Validate())
}

func TestValidateMissingServiceName(t *testing.T) {
	cfg := BaseConfig{
		ServiceName: "",
		DatabaseURL: "postgres://localhost/db",
		HTTPPort:    8080,
	}
	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "SERVICE_NAME")
}

func TestValidateMissingDatabaseURL(t *testing.T) {
	cfg := BaseConfig{
		ServiceName: "my-service",
		DatabaseURL: "",
		HTTPPort:    8080,
	}
	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "DATABASE_URL")
}

func TestValidateInvalidPortZero(t *testing.T) {
	cfg := BaseConfig{
		ServiceName: "my-service",
		DatabaseURL: "postgres://localhost/db",
		HTTPPort:    0,
	}
	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "HTTP_PORT")
}

func TestValidateInvalidPortNegative(t *testing.T) {
	cfg := BaseConfig{
		ServiceName: "my-service",
		DatabaseURL: "postgres://localhost/db",
		HTTPPort:    -1,
	}
	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "HTTP_PORT")
}

func TestValidateInvalidPortTooHigh(t *testing.T) {
	cfg := BaseConfig{
		ServiceName: "my-service",
		DatabaseURL: "postgres://localhost/db",
		HTTPPort:    70000,
	}
	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "HTTP_PORT")
}

func TestValidatePortBoundaries(t *testing.T) {
	// Port 1 — valid minimum
	cfg := BaseConfig{
		ServiceName: "svc",
		DatabaseURL: "postgres://localhost/db",
		HTTPPort:    1,
	}
	assert.NoError(t, cfg.Validate())

	// Port 65535 — valid maximum
	cfg.HTTPPort = 65535
	assert.NoError(t, cfg.Validate())

	// Port 65536 — one past max
	cfg.HTTPPort = 65536
	assert.Error(t, cfg.Validate())
}

func TestBaseConfigGetters(t *testing.T) {
	cfg := BaseConfig{
		ServiceName:      "svc",
		OTelEnabled:      true,
		OTelCollectorURL: "endpoint:4317",
		DatabaseURL:      "postgres://localhost/db",
		KafkaBrokers:     "broker:9092",
	}
	assert.Equal(t, "svc", cfg.GetServiceName())
	assert.True(t, cfg.GetOTelEnabled())
	assert.Equal(t, "endpoint:4317", cfg.GetOTelEndpoint())
	assert.Equal(t, "postgres://localhost/db", cfg.GetDatabaseURL())
	assert.Equal(t, "broker:9092", cfg.GetKafkaBrokers())
}

func TestLoadBaseConfigFromReader(t *testing.T) {
	yamlContent := `
SERVICE_NAME: "reader-test-svc"
ENV: "staging"
HTTP_PORT: 8888
LOG_LEVEL: "warn"
`
	cfg, err := LoadBaseConfigFromReader(strings.NewReader(yamlContent))
	require.NoError(t, err)
	assert.Equal(t, "reader-test-svc", cfg.ServiceName)
	assert.Equal(t, "staging", cfg.Environment)
	assert.Equal(t, 8888, cfg.HTTPPort)
	assert.Equal(t, "warn", cfg.LogLevel)
}
