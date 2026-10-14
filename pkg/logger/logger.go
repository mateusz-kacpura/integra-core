package logger

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/mateusz-kacpura/integra-core/pkg/telemetry"
)

// EntrypointConfig defines operational and bootstrap metadata for structured logging.
type EntrypointConfig struct {
	Service     string
	Component   string
	ProcessType string // "server", "worker", "migrate", "admin", etc.
	Environment string
	Version     string
	LogLevel    string
}

// NewBootstrapLogger creates and registers a zerolog.Logger configured with standard
// operational metadata: service name (app.name), process type (component), environment, and version,
// along with OpenTelemetry trace and PII sanitizer hooks.
func NewBootstrapLogger(cfg EntrypointConfig) zerolog.Logger {
	_, l := InitEntrypoint(context.Background(), cfg)
	return l
}

// InitEntrypoint initializes Zerolog with full operational context (service, component,
// process_type, environment, version, correlation ID), attaches OpenTelemetry trace and PII sanitizer hooks,
// updates global log.Logger, and returns the contextualized context and logger writing to os.Stdout.
func InitEntrypoint(ctx context.Context, cfg EntrypointConfig) (context.Context, zerolog.Logger) {
	return InitEntrypointWithWriter(ctx, cfg, os.Stdout)
}

// InitEntrypointWithWriter initializes Zerolog with full operational context using the provided writer.
func InitEntrypointWithWriter(ctx context.Context, cfg EntrypointConfig, out io.Writer) (context.Context, zerolog.Logger) {
	if ctx == nil {
		ctx = context.Background()
	}
	zerolog.TimeFieldFormat = time.RFC3339Nano

	env := cfg.Environment
	if env == "" {
		env = os.Getenv("INTEGRA_ENV")
	}
	if env == "" {
		env = os.Getenv("ENV")
	}
	if env == "" {
		env = "development"
	}

	service := cfg.Service
	if service == "" {
		service = os.Getenv("INTEGRA_SERVICE_NAME")
	}
	if service == "" {
		service = os.Getenv("SERVICE_NAME")
	}

	processType := cfg.ProcessType
	if processType == "" {
		processType = cfg.Component
	}
	if processType == "" {
		processType = os.Getenv("INTEGRA_PROCESS_TYPE")
	}
	if processType == "" {
		processType = os.Getenv("INTEGRA_COMPONENT")
	}

	version := cfg.Version
	if version == "" {
		version = os.Getenv("INTEGRA_VERSION")
	}
	if version == "" {
		version = "v0.2.0-dev.0"
	}

	logLevelStr := cfg.LogLevel
	if logLevelStr == "" {
		logLevelStr = os.Getenv("INTEGRA_LOG_LEVEL")
	}
	if logLevelStr == "" {
		logLevelStr = os.Getenv("LOG_LEVEL")
	}
	level, err := zerolog.ParseLevel(logLevelStr)
	if err != nil || logLevelStr == "" {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	var writer io.Writer = out
	if env == "development" && os.Getenv("INTEGRA_LOG_FORMAT") == "console" {
		writer = zerolog.ConsoleWriter{Out: out, TimeFormat: time.RFC3339}
	}
	writer = NewSanitizingWriter(writer)

	// Ensure a correlation ID exists on context
	corrID := telemetry.CorrelationID(ctx)
	if corrID == "" {
		ctx, corrID = telemetry.WithCorrelationID(ctx, "")
	}

	builder := zerolog.New(writer).Level(level).With().Timestamp()
	if service != "" {
		builder = builder.Str("service", service).Str("app.name", service)
	}
	if processType != "" {
		builder = builder.Str("process_type", processType).Str("component", processType)
	}
	if env != "" {
		builder = builder.Str("environment", env)
	}
	if version != "" {
		builder = builder.Str("version", version)
	}
	if corrID != "" {
		builder = builder.Str("correlation_id", corrID)
	}

	l := builder.Logger().
		Hook(telemetry.NewTraceHook()).
		Hook(SanitizerHook{})

	log.Logger = l
	ctx = l.WithContext(ctx)
	return ctx, l
}

// InitLogger creates a zerolog logger with integra defaults (zero-allocation, JSON format).
func InitLogger(level string, isLocalDev bool) zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339Nano

	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(logLevel)

	var zl zerolog.Logger
	if isLocalDev {
		zl = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).With().Timestamp().Logger()
	} else {
		// Raw JSON logging for collection by Loki/FluentBit
		zl = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}
	return zl
}

// NewWithTraceHook creates a zerolog.Logger with the OTel TraceHook attached.
// The TraceHook automatically injects trace_id and span_id from context into log entries.
func NewWithTraceHook(serviceName string) zerolog.Logger {
	return zerolog.New(os.Stdout).With().
		Timestamp().
		Str("service", serviceName).
		Logger().
		Hook(telemetry.NewTraceHook())
}

// NewWithSanitizer creates a zerolog.Logger with both OTel TraceHook and
// SanitizerHook attached. The SanitizerHook detects PII in log messages and
// appends a "sanitized_msg" field when PII is found.
//
// For complete PII protection, wrap untrusted input with SanitizeString before
// passing it to log methods:
//
//	log.Info().Msg(logger.SanitizeString(untrustedInput))
func NewWithSanitizer(serviceName string) zerolog.Logger {
	writer := NewSanitizingWriter(os.Stdout)
	return zerolog.New(writer).With().
		Timestamp().
		Str("service", serviceName).
		Logger().
		Hook(telemetry.NewTraceHook()).
		Hook(SanitizerHook{})
}
