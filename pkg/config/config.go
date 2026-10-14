package config

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// DefaultShutdownGrace is the default duration to sleep after cancelling
// a context, giving in-flight work time to complete before process exit.
const DefaultShutdownGrace = 2 * time.Second

// BaseConfig contains the fields shared across all integra microservices.
// Services embed this struct into their own config and add service-specific
// fields on top. BaseConfig implements fxutil.ServiceConfig.
type BaseConfig struct {
	ServiceName      string `mapstructure:"SERVICE_NAME"`
	Environment      string `mapstructure:"ENV"`
	HTTPPort         int    `mapstructure:"HTTP_PORT"`
	GRPCPort         int    `mapstructure:"GRPC_PORT"`
	LogLevel         string `mapstructure:"LOG_LEVEL"`
	DatabaseURL      string `mapstructure:"DATABASE_URL"`
	KafkaBrokers     string `mapstructure:"KAFKA_BROKERS"`
	ValkeyURL        string `mapstructure:"VALKEY_URL"`
	OTelEnabled      bool   `mapstructure:"OTEL_ENABLED"`
	OTelCollectorURL string `mapstructure:"OTEL_COLLECTOR_URL"`
}

// GetServiceName returns the configured microservice name.
func (c BaseConfig) GetServiceName() string { return c.ServiceName }

// GetOTelEnabled reports whether OpenTelemetry instrumentation is enabled.
func (c BaseConfig) GetOTelEnabled() bool { return c.OTelEnabled }

// GetOTelEndpoint returns the configured OpenTelemetry collector URL endpoint.
func (c BaseConfig) GetOTelEndpoint() string { return c.OTelCollectorURL }

// GetDatabaseURL returns the connection URL string for the primary database.
func (c BaseConfig) GetDatabaseURL() string { return c.DatabaseURL }

// GetKafkaBrokers returns the comma-separated Kafka bootstrap broker addresses.
func (c BaseConfig) GetKafkaBrokers() string { return c.KafkaBrokers }

// Validate checks that required BaseConfig fields are present and within
// acceptable ranges. Call this after loading the config.
func (c BaseConfig) Validate() error {
	if c.ServiceName == "" {
		return fmt.Errorf("SERVICE_NAME is required")
	}
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	if c.HTTPPort <= 0 || c.HTTPPort > 65535 {
		return fmt.Errorf("HTTP_PORT must be between 1 and 65535, got %d", c.HTTPPort)
	}
	return nil
}

// AppConfig extends BaseConfig with the original convenience fields (AppName,
// Port) so that existing callers of LoadConfig continue to compile.
//
// Deprecated: New services should embed BaseConfig directly in their own
// config struct and use SERVICE_NAME / HTTP_PORT instead of APP_NAME / PORT.
type AppConfig struct {
	BaseConfig `mapstructure:",squash"`

	// Legacy aliases kept for backward compatibility.
	AppName string `mapstructure:"APP_NAME"`
	Port    int    `mapstructure:"PORT"`
}

// LoadConfig loads the integra system configuration from environment variables.
func LoadConfig() (*AppConfig, error) {
	return LoadConfigFromReader(nil)
}

// LoadConfigFromReader loads AppConfig from an optional io.Reader stream and environment variables.
func LoadConfigFromReader(r io.Reader) (*AppConfig, error) {
	v := viper.New()
	v.SetEnvPrefix("integra")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Defaults for the legacy AppConfig path.
	v.SetDefault("APP_NAME", "integra-service")
	v.SetDefault("ENV", "development")
	v.SetDefault("PORT", 8080)
	v.SetDefault("HTTP_PORT", 8080)
	v.SetDefault("LOG_LEVEL", "info")

	if r != nil {
		v.SetConfigType("yaml")
		if err := v.ReadConfig(r); err != nil {
			return nil, fmt.Errorf("read config from reader: %w", err)
		}
	}

	var config AppConfig
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	// If the new ServiceName was not set explicitly, fall back to AppName.
	if config.ServiceName == "" {
		config.ServiceName = config.AppName
	}
	// If HTTPPort was not set explicitly, fall back to the legacy Port.
	if config.HTTPPort == 0 && config.Port != 0 {
		config.HTTPPort = config.Port
	}

	return &config, nil
}

// LoadBaseConfig loads only the BaseConfig fields. This is the recommended
// entry-point for new services that embed BaseConfig directly.
func LoadBaseConfig() (*BaseConfig, error) {
	return LoadBaseConfigFromReader(nil)
}

// LoadBaseConfigFromReader loads BaseConfig from an io.Reader stream and environment variables.
func LoadBaseConfigFromReader(r io.Reader) (*BaseConfig, error) {
	v := viper.New()
	v.SetEnvPrefix("integra")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.SetDefault("SERVICE_NAME", "integra-service")
	v.SetDefault("ENV", "development")
	v.SetDefault("HTTP_PORT", 8080)
	v.SetDefault("LOG_LEVEL", "info")

	if r != nil {
		v.SetConfigType("yaml")
		if err := v.ReadConfig(r); err != nil {
			return nil, fmt.Errorf("read config from reader: %w", err)
		}
	}

	var config BaseConfig
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
