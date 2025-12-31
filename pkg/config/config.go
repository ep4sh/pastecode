package config

import (
	"fmt"

	"github.com/caarlos0/env"
)

// BackendConfig defines web server configuration.
type BackendConfig struct {
	Host string `env:"PASTECODE_HOST" envDefault:"0.0.0.0"`
	Port string `env:"PASTECODE_PORT" envDefault:"8080"`
}

// NewBackendConfig is a web server config fabric.
func NewBackendConfig() (*BackendConfig, error) {
	var c BackendConfig
	if err := env.Parse(&c); err != nil {
		return nil, err
	}

	return &c, nil
}

// Addr returns BackendConfig addr string for serving http requests.
func (c BackendConfig) Addr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
