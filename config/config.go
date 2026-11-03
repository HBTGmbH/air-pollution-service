// Package config reads the service configuration from environment variables.
package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

// Conf is the configuration of the service.
type Conf struct {
	AirPollutionFile string `env:"AIR_POLLUTION_FILE" envDefault:"air-pollution.csv"`
	Server           ConfServer
}

// ConfServer is the configuration of the HTTP server.
type ConfServer struct {
	Port                  int           `env:"SERVER_PORT"                    envDefault:"8080"`
	IdleTimeout           time.Duration `env:"SERVER_IDLE_TIMEOUT"            envDefault:"60s"`
	ReadHeaderTimeout     time.Duration `env:"SERVER_READ_HEADER_TIMEOUT"     envDefault:"10s"`
	ShutdownSleepDuration time.Duration `env:"SERVER_SHUTDOWN_SLEEP_DURATION" envDefault:"0s"`
	DrainTimeout          time.Duration `env:"SERVER_DRAIN_TIMEOUT"           envDefault:"30s"`
	ShutdownTimeout       time.Duration `env:"SERVER_SHUTDOWN_TIMEOUT"        envDefault:"30s"`
}

// New reads the configuration from the environment.
func New() (*Conf, error) {
	cfg, err := env.ParseAs[Conf]()
	if err != nil {
		return nil, fmt.Errorf("parsing environment variables: %w", err)
	}

	return &cfg, nil
}
