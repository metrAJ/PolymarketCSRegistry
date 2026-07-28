package config

import (
	"os"
)

type Config struct {
	Port string
}

func LoadConfig() (*Config, error) {
	return &Config{
		Port: os.Getenv("PORT"),
	}, nil
}
