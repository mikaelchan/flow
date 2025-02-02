package main

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Server struct {
		BaseURL string `toml:"base_url"`
		Timeout int    `toml:"timeout"`
	}
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	if err := toml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
