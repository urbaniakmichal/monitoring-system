//go:build unit

package config

import (
	"testing"
	"time"
)

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
{
			name: "Valid configuration",
			cfg: Config{
				Interval: 5 * time.Second,
				Timeout:  2 * time.Second,
				Server: ServerConfig{
					Port:         8080,
					ReadTimeout:  5 * time.Second,
					WriteTimeout: 10 * time.Second,
					IdleTimeout:  60 * time.Second,
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid interval (zero or negative)",
			cfg: Config{
				Interval: 0,
				Timeout:  2 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "Invalid timeout greater than interval",
			cfg: Config{
				Interval: 2 * time.Second,
				Timeout:  5 * time.Second,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}