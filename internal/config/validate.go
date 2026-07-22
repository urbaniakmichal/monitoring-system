package config

import "errors"

func (c Config) Validate() error {

	if c.Interval <= 0 {
		return errors.New("interval must be greater than zero")
	}

	if c.Timeout <= 0 {
		return errors.New("timeout must be greater than zero")
	}

	return nil
}
