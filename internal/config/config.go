package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"
)

// ServerConfig holds configuration values for the server.
type ServerConfig struct {
	Address      string `env:"ADDRESS"`                                                                                              // Address to listen on (e.g., ":8080")
	DatabaseURI  string `env:"DATABASE_URI" envDefault:"postgres://gophkeeper:gophkeeper@localhost:5432/gophkeeper?sslmode=disable"` // URI for the database connection
	LogLevel     string `env:"LOG_LEVEL"`                                                                                            // Log level (e.g., "info")
	JWTSecretKey string `env:"JWT_SECRET_KEY"`                                                                                       // JWT secret key used for authentication
}

// NewServerConfig creates and returns a new instance of ServerConfig.
func NewServerConfig() *ServerConfig {
	return &ServerConfig{}
}

// GetConfigFromFile loads configuration from a JSON file.
//
// Parameters:
//   - configPath (string): The path to the JSON configuration file.
//
// Returns:
//   - error: An error if the file cannot be read or if the JSON cannot be unmarshalled.
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

// Init initializes the ServerConfig by parsing command-line flags and environment variables.
// It also loads configuration from a file if the CONFIG environment variable is set.
//
// Returns:
//   - *ServerConfig: The initialized ServerConfig instance.
//   - error: An error if any operation fails during initialization.
func (c *ServerConfig) Init() (*ServerConfig, error) {
	// Define command-line flags with default values
	flag.StringVar(&c.Address, "a", ":8080", "address to listen on")
	flag.StringVar(&c.DatabaseURI, "d", "", "database URI")
	flag.StringVar(&c.LogLevel, "l", "info", "log level")
	flag.StringVar(&c.JWTSecretKey, "j", "", "auth secret key")
	flag.Parse()

	// Check if a configuration file path is provided in the CONFIG environment variable
	if config := os.Getenv("CONFIG"); config != "" {
		err := c.GetConfigFromFile(config)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file '%s': %w", config, err)
		}
	}

	// Parse environment variables and override command-line flag values
	if err := env.Parse(c); err != nil {
		return nil, fmt.Errorf("failed to parse env: %w", err)
	}

	return c, nil
}
