package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env"
)

type ServerConfig struct {
	Address      string `env:"ADDRESS"`
	DatabaseURI  string `env:"DATABASE_URI" envDefault:"postgres://dev-backend:dev-backend@postgres:5432/dev-backend?sslmode=disable"`
	LogLevel     string `env:"LOG_LEVEL"`
	JWTSecretKey string `env:"JWT_SECRET_KEY"`
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{}
}

func (c *ServerConfig) GetConfigFromFile(configPath string) error {
	b, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read file by path '%s': %w", configPath, err)
	}

	err = json.Unmarshal(b, c)
	if err != nil {
		return fmt.Errorf("failed to unmarshal cfg '%s': %w", string(b), err)
	}

	return nil
}

func (c *ServerConfig) Init() (*ServerConfig, error) {
	flag.StringVar(&c.Address, "a", ":8080", "address to listen on")
	flag.StringVar(&c.DatabaseURI, "d", "", "database URI")
	flag.StringVar(&c.LogLevel, "l", "info", "log level")
	flag.StringVar(&c.JWTSecretKey, "j", "", "auth secret key")
	flag.Parse()

	if err := env.Parse(c); err != nil {
		return nil, fmt.Errorf("failed to parse env: %w", err)
	}

	return c, nil
}
