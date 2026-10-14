// Package baseconfig re-exports common base configuration structures and loaders
// to avoid qualifier collisions when importing alongside service-local config packages.
package baseconfig

import (
	"io"

	"github.com/mateusz-kacpura/integra-core/pkg/config"
)

// DefaultShutdownGrace re-exports config.DefaultShutdownGrace.
const DefaultShutdownGrace = config.DefaultShutdownGrace

// BaseConfig re-exports config.BaseConfig.
type BaseConfig = config.BaseConfig

// LoadBaseConfig re-exports config.LoadBaseConfig.
func LoadBaseConfig() (*BaseConfig, error) {
	return config.LoadBaseConfig()
}

// LoadBaseConfigFromReader re-exports config.LoadBaseConfigFromReader.
func LoadBaseConfigFromReader(r io.Reader) (*BaseConfig, error) {
	return config.LoadBaseConfigFromReader(r)
}
