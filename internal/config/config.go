package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

func Load() Config {
	cfg := Config{}

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	return cfg
}
