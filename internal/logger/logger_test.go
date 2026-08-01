//go:build unit

package logger

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"testing"
)

func TestNewLoggerJSON(t *testing.T) {
	var buf bytes.Buffer

	cfg := Config{
		Level:  "info",
		Format: "json",
	}

	log := NewLogger(cfg, &buf)
	log.Info("application_started", slog.String("version", "1.0.0"), slog.Int("port", 8080))

	output := buf.String()

	var entry map[string]any
	if err := json.Unmarshal([]byte(output), &entry); err != nil {
		t.Fatalf("failed to parse log output as JSON: %v. Output was: %s", err, output)
	}

	if entry["msg"] != "application_started" {
		t.Errorf("expected message 'application_started', got '%v'", entry["msg"])
	}

	if entry["version"] != "1.0.0" {
		t.Errorf("expected version '1.0.0', got '%v'", entry["version"])
	}

	if code, ok := entry["port"].(float64); !ok || code != 8080 {
		t.Errorf("expected port 8080, got '%v'", entry["port"])
	}
}

func TestNewLoggerLevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	cfg := Config{
		Level:  "warn",
		Format: "json",
	}

	log := NewLogger(cfg, &buf)

	log.Info("debug info message")

	log.Warn("warning alert triggered")

	output := buf.String()
	if buf.Len() == 0 {
		t.Fatal("expected output for warn level, got empty buffer")
	}

	var entry map[string]any
	if err := json.Unmarshal([]byte(output), &entry); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if entry["msg"] != "warning alert triggered" {
		t.Errorf("expected warning message, got '%v'", entry["msg"])
	}
}
