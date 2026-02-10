package config

import (
	"log"
	"time"

	"github.com/caarlos0/env/v11"
)

type Conf struct {
	Server           ConfServer
	AirPollutionFile string `env:"AIR_POLLUTION_FILE" envDefault:"air-pollution.csv"`
}

type ConfServer struct {
	Port                  int           `env:"SERVER_PORT" envDefault:"8080"`
	IdleTimeout           time.Duration `env:"SERVER_IDLE_TIMEOUT" envDefault:"60s"`
	ShutdownSleepDuration time.Duration `env:"SERVER_SHUTDOWN_SLEEP_DURATION" envDefault:"0s"`
	DrainTimeout          time.Duration `env:"SERVER_DRAIN_TIMEOUT" envDefault:"30s"`
	ShutdownTimeout       time.Duration `env:"SERVER_SHUTDOWN_TIMEOUT" envDefault:"30s"`
}

func New() *Conf {
	var cfg Conf
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatalf("Failed to parse environment variables: %s", err)
	}
	return &cfg
}
