package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)


type Config struct {
	
	PingInterval time.Duration `json:"ping_interval_ms"`

	Concurrency int `json:"concurrency"`

	// RequestTimeout is the maximum time to wait for a service response.
	RequestTimeout time.Duration `json:"request_timeout_ms"`

	// Services is the list of URLs to monitor.
	Services []string `json:"services"`
}

// LoadConfiguration reads a JSON configuration file and returns a Config struct.
// It uses os.Getenv as a simple layering mechanism to override file values.
func LoadConfiguration(file string) (*Config, error) {
	var config Config

	configFile, err := os.Open(file)
	if err!= nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err!= nil {
		return nil, fmt.Errorf("failed to decode JSON config: %w", err)
	}

	// Override JSON values with environment variables if they exist.
	// This provides a critical layered configuration approach for production.
	if val, ok := os.LookupEnv("PING_INTERVAL_MS"); ok {
		if d, err := time.ParseDuration(val + "ms"); err == nil {
			config.PingInterval = d
		}
	}
	if val, ok := os.LookupEnv("CONCURRENCY"); ok {
		if n, err := fmt.Sscan(val, &config.Concurrency); err == nil && n == 1 {
			// Sscan returns the number of items successfully scanned.
		}
	}
	if val, ok := os.LookupEnv("REQUEST_TIMEOUT_MS"); ok {
		if d, err := time.ParseDuration(val + "ms"); err == nil {
			config.RequestTimeout = d
		}
	}

	// Convert from milliseconds for readability in JSON.
	config.PingInterval *= time.Millisecond
	config.RequestTimeout *= time.Millisecond

	return &config, nil
}