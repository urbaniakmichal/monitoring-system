package config

import (
	"errors"
	"fmt"
)

func (c Config) Validate() error {

	if c.Interval <= 0 {
		return errors.New("interval must be greater than zero")
	}

	if c.Timeout <= 0 {
		return errors.New("timeout must be greater than zero")
	}

	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d (must be between 1 and 65535)", c.Server.Port)
	}

	if c.Server.ReadTimeout <= 0 {
		return fmt.Errorf("invalid read_timeout: %v", c.Server.ReadTimeout)
	}

	if c.Server.WriteTimeout <= 0 {
		return fmt.Errorf("invalid write_timeout: %v", c.Server.WriteTimeout)
	}

	if c.Server.IdleTimeout <= 0 {
		return fmt.Errorf("invalid idle_timeout: %v", c.Server.IdleTimeout)
	}

	return nil
}
