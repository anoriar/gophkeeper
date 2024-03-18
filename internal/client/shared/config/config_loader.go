package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

// LoadConfig missing godoc.
func LoadConfig() (*Config, error) {
	conf := NewConfig()

	err := env.Parse(conf)
	if err != nil {
		return nil, fmt.Errorf("parse env error: %v", err)
	}

	return conf, nil
}
