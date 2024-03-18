package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

// LoadConfig missing godoc.
func LoadConfig() (*Config, error) {
	conf := NewConfig()
	parseFlags(conf)

	err := env.Parse(conf)
	if err != nil {
		return nil, fmt.Errorf("parse env error: %v", err)
	}

	return conf, nil
}

func parseFlags(config *Config) {
	flag.StringVar(&config.RunAddress, "a", "localhost:8080", "RunAddress")
	flag.StringVar(&config.DatabaseURI, "d", "", "Database DSN")
	flag.StringVar(&config.JwtSecretKey, "j", "secret-key", "Auth secret key")

	flag.Parse()
}
