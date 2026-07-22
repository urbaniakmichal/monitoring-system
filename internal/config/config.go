package config

import (
	"os"
	"path/filepath"
	"time"

	"go.yaml.in/yaml/v3"
)

type Config struct {
	Interval time.Duration `yaml:"interval"`
	Timeout  time.Duration `yaml:"timeout"`
}

func DefaultConfig() Config {
	return Config{
		Interval: 5 * time.Second,
		Timeout:  2 * time.Second,
	}
}

func Load(path string) (*Config, error) {
	cleanPath := filepath.Clean(path)

	data, err := os.ReadFile(cleanPath)
	if err != nil {
		return nil, err
	}

	var config Config

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
