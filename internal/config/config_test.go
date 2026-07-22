package config

import "testing"

func TestLoad(t *testing.T) {

	cfg, err := Load("../../configs/config.yaml")

	if err != nil {
		t.Fatal(err)
	}

	if cfg.Interval.String() != "5s" {
		t.Fatalf("wrong interval: %v", cfg.Interval)
	}

	if cfg.Timeout.String() != "2s" {
		t.Fatalf("wrong timeout: %v", cfg.Timeout)
	}
}
