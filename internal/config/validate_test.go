package config

import (
	"testing"
	"time"
)

func TestConfigValidate_OK(t *testing.T) {

	config := Config{
		Interval: 5 * time.Second,
		Timeout:  2 * time.Second,
	}

	err := config.Validate()

	if err != nil {
		t.Fatal(err)
	}
}

func TestConfigValidate_InvalidInterval(t *testing.T) {

	config := Config{
		Interval: 0,
		Timeout:  2 * time.Second,
	}

	err := config.Validate()

	if err == nil {
		t.Fatal("expected error")
	}
}
