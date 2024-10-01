package config

import (
	"github.com/joeshaw/envdecode"
	"log"
)

type Conf struct {
	Server                      ConfServer
	AirPollutionFile            string `env:"AIR_POLLUTION_FILE,default=air-pollution.csv"`
	SleepDurationBeforeShutdown string `env:"SLEEP_DURATION_BEFORE_SHUTDOWN,default=10s"`
}

type ConfServer struct {
	Port int `env:"SERVER_PORT,default=8080"`
}

func New() *Conf {
	var c Conf
	if err := envdecode.StrictDecode(&c); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}

	return &c
}
