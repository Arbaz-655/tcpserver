package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ServerConfig ServerConfig `yaml:"server"`
}

type ServerConfig struct {
	Ports []int `yaml:"ports"`
}

func Load(filename string) (*Config, error) {
	// Read the YAML file
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	// Create a Config struct to hold the parsed data
	var config Config

	// Unmarshal the YAML into the Config struct
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}

	// Validate the configuration
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %v", err)
	}

	return &config, nil
}

func validateConfig(config *Config) error {
	// Validate Server configuration
	if len(config.ServerConfig.Ports) == 0 {
		return fmt.Errorf("at least one server port is required")
	}

	return nil
}
