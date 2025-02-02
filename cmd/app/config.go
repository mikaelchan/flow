package main

import (
	"os"

	"github.com/mikaelchan/hamster/internal/application"
	"github.com/pelletier/go-toml/v2"
)

type AppConfig struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

type Config struct {
	App     AppConfig                 `toml:"app"`
	Service application.ServiceConfig `toml:"service"`
}

func LoadConfig(path string) (Config, error) {
	var config Config
	data, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}
	err = toml.Unmarshal(data, &config)
	return config, err
}
