package config

import (
	"fmt"

	"github.com/caarlos0/env/v8"
)

type Config struct {
	Port int `env:"PORT" envDefault:"10000"`
}

func NewConfig() (*Config, error) {
	var c Config
	if err := env.Parse(&c); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return &c, nil
}
