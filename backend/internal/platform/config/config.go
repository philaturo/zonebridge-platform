// Package config provides configuration loading and validation for the ZoneBridge platform.
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	platformerrors "github.com/philaturo/zonebridge-platform/internal/platform/errors"
)

// Config holds all configuration values for the platform.
type Config struct {
	ServerHost            string
	ServerPort            int
	ServerReadTimeout     time.Duration
	ServerWriteTimeout    time.Duration
	ServerShutdownTimeout time.Duration
	LogLevel              string
	LogFormat             string
	Environment           string
}

// Load reads configuration from environment variables and returns a validated Config.
func Load() (*Config, error) {
	cfg := &Config{
		ServerHost:  getEnv("ZONEBRIDGE_SERVER_HOST", ""),
		LogLevel:    getEnv("ZONEBRIDGE_LOG_LEVEL", "info"),
		LogFormat:   getEnv("ZONEBRIDGE_LOG_FORMAT", "json"),
		Environment: getEnv("ZONEBRIDGE_ENVIRONMENT", "production"),
	}

	var err error

	cfg.ServerPort, err = getEnvInt("ZONEBRIDGE_SERVER_PORT", 8080)
	if err != nil {
		return nil, fmt.Errorf("load server port: %w", err)
	}

	cfg.ServerReadTimeout, err = getEnvDuration("ZONEBRIDGE_SERVER_READ_TIMEOUT", 10*time.Second)
	if err != nil {
		return nil, fmt.Errorf("load server read timeout: %w", err)
	}

	cfg.ServerWriteTimeout, err = getEnvDuration("ZONEBRIDGE_SERVER_WRITE_TIMEOUT", 10*time.Second)
	if err != nil {
		return nil, fmt.Errorf("load server write timeout: %w", err)
	}

	cfg.ServerShutdownTimeout, err = getEnvDuration("ZONEBRIDGE_SERVER_SHUTDOWN_TIMEOUT", 30*time.Second)
	if err != nil {
		return nil, fmt.Errorf("load server shutdown timeout: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("validate configuration: %w", err)
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.ServerPort < 1 || c.ServerPort > 65535 {
		return fmt.Errorf("%w: server port must be between 1 and 65535, got %d", platformerrors.ErrInvalidConfiguration, c.ServerPort)
	}
	if c.ServerReadTimeout <= 0 {
		return fmt.Errorf("%w: server read timeout must be positive", platformerrors.ErrInvalidConfiguration)
	}
	if c.ServerWriteTimeout <= 0 {
		return fmt.Errorf("%w: server write timeout must be positive", platformerrors.ErrInvalidConfiguration)
	}
	if c.ServerShutdownTimeout <= 0 {
		return fmt.Errorf("%w: server shutdown timeout must be positive", platformerrors.ErrInvalidConfiguration)
	}

	validLogLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLogLevels[c.LogLevel] {
		return fmt.Errorf("%w: log level must be one of [debug, info, warn, error], got %q", platformerrors.ErrInvalidConfiguration, c.LogLevel)
	}

	validLogFormats := map[string]bool{"json": true, "text": true}
	if !validLogFormats[c.LogFormat] {
		return fmt.Errorf("%w: log format must be one of [json, text], got %q", platformerrors.ErrInvalidConfiguration, c.LogFormat)
	}

	validEnvironments := map[string]bool{"development": true, "production": true, "test": true}
	if !validEnvironments[c.Environment] {
		return fmt.Errorf("%w: environment must be one of [development, production, test], got %q", platformerrors.ErrInvalidConfiguration, c.Environment)
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) (int, error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue, nil
	}
	
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%w: environment variable %q must be an integer, got %q", platformerrors.ErrInvalidConfiguration, key, value)
	}
	
	return intValue, nil
}

func getEnvDuration(key string, defaultValue time.Duration) (time.Duration, error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue, nil
	}
	
	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("%w: environment variable %q must be a valid duration, got %q", platformerrors.ErrInvalidConfiguration, key, value)
	}
	
	return duration, nil
}