package logger

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"time"
)

type LogEntry struct {
	Timestamp time.Time `json:"@timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Service   string    `json:"service"`
}

func SendToElastic(level, message string) {
	esURL := os.Getenv("ELASTICSEARCH_URL")
	if esURL == "" {
		esURL = "http://localhost:9200"
	}

	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
		Service:   "monitor-system",
	}

	payload, err := json.Marshal(entry)
	if err != nil {
		return
	}

	// index: monitor-system-logs
	go func() {
		client := &http.Client{Timeout: 2 * time.Second}
		_, _ = client.Post(esURL+"/monitor-system-logs/_doc", "application/json", bytes.NewBuffer(payload))
	}()
}
