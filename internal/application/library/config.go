package library

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

type PostgresConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	Database string `toml:"database"`
}

type RedisConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
}

type Config struct {
	Postgres PostgresConfig `toml:"postgres"`
	Redis    RedisConfig    `toml:"redis"`
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
