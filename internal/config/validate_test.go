package config

import (
	"testing"
	"time"
)

func validConfig() Config {
	return Config{
		Interval: 5 * time.Second,
		Timeout:  2 * time.Second,
		Server: ServerConfig{
			Port:         8080,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
}

func TestConfigValidate_OK(t *testing.T) {
	cfg := validConfig()

	err := cfg.Validate()
	if err != nil {
		t.Fatalf("expected valid config, got error: %v", err)
	}
}

func TestConfigValidate_InvalidInterval(t *testing.T) {
	cfg := validConfig()
	cfg.Interval = 0

	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for invalid interval, got nil")
	}
}

func TestConfigValidate_InvalidServerPort(t *testing.T) {
	testCases := []struct {
		name string
		port int
	}{
		{"Port zero", 0},
		{"Port negative", -1},
		{"Port too high", 70000},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := validConfig()
			cfg.Server.Port = tc.port

			err := cfg.Validate()
			if err == nil {
				t.Fatalf("expected error for port %d, got nil", tc.port)
			}
		})
	}
}

func TestConfigValidate_InvalidServerTimeouts(t *testing.T) {
	t.Run("Zero ReadTimeout", func(t *testing.T) {
		cfg := validConfig()
		cfg.Server.ReadTimeout = 0

		if err := cfg.Validate(); err == nil {
			t.Fatal("expected error for zero read_timeout, got nil")
		}
	})

	t.Run("Zero WriteTimeout", func(t *testing.T) {
		cfg := validConfig()
		cfg.Server.WriteTimeout = 0

		if err := cfg.Validate(); err == nil {
			t.Fatal("expected error for zero write_timeout, got nil")
		}
	})

	t.Run("Zero IdleTimeout", func(t *testing.T) {
		cfg := validConfig()
		cfg.Server.IdleTimeout = 0

		if err := cfg.Validate(); err == nil {
			t.Fatal("expected error for zero idle_timeout, got nil")
		}
	})
}
