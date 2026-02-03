package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Redis       RedisConfig               `yaml:"redis"`
	Mappings    map[string][]string       `yaml:"-"` // Normalized mappings (always []string)
	RawMappings map[string]interface{}    `yaml:"mappings"` // Raw mappings from YAML
}

// RedisConfig contains Redis connection settings
type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// LoadConfig loads configuration from a YAML file
func LoadConfig(filepath string) (*Config, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if config.Redis.Password == "" {
		config.Redis.Password = os.Getenv("REDIS_PASSWORD")
	}

	// Normalize mappings to always be []string
	if err := normalizeMappings(&config); err != nil {
		return nil, fmt.Errorf("failed to normalize mappings: %w", err)
	}

	return &config, nil
}

// normalizeMappings converts raw mappings to normalized []string format
// Supports both single string values and array of strings
func normalizeMappings(config *Config) error {
	config.Mappings = make(map[string][]string)
	
	for channel, value := range config.RawMappings {
		switch v := value.(type) {
		case string:
			// Single queue (backward compatibility)
			config.Mappings[channel] = []string{v}
		case []interface{}:
			// Multiple queues
			queues := make([]string, 0, len(v))
			for i, item := range v {
				queueStr, ok := item.(string)
				if !ok {
					return fmt.Errorf("mapping for channel '%s' contains non-string value at index %d", channel, i)
				}
				queues = append(queues, queueStr)
			}
			if len(queues) == 0 {
				return fmt.Errorf("mapping for channel '%s' has empty queue list", channel)
			}
			config.Mappings[channel] = queues
		default:
			return fmt.Errorf("mapping for channel '%s' has invalid type: expected string or array", channel)
		}
	}
	
	return nil
}
