package runner

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"

	"monitoring-system/internal/config"
	"monitoring-system/internal/memory_storage"
	"monitoring-system/internal/metrics"
)

type contextKey string

const TraceIDKey contextKey = "trace_id"

type Runner struct {
	Config  config.Config
	Logger  *slog.Logger
	Storage *memory_storage.MemoryStorage
}

func (r *Runner) Start(ctx context.Context) {
	logger := r.Logger
	if logger == nil {
		logger = slog.Default()
	}

	if r.Config.Interval <= 0 {
		r.Config.Interval = 10 * time.Second
	}

	if r.Config.Timeout <= 0 {
		r.Config.Timeout = 5 * time.Second
	}

	ticker := time.NewTicker(r.Config.Interval)
	defer ticker.Stop()

	logger.InfoContext(ctx, "starting metrics collector", slog.Duration("interval", r.Config.Interval))
	r.collect(ctx, logger) // Initial collection before ticker triggers

	for {
		select {
		case <-ctx.Done():
			logger.InfoContext(ctx, "stopping metrics collector")
			return
		case <-ticker.C:
			r.collect(ctx, logger)
		}
	}
}

func (r *Runner) collect(ctx context.Context, logger *slog.Logger) {
	// 1. Generate Trace ID before starting collection
	traceID := generateTraceID()

	// 2. Attach Trace ID to context and create a logger instance with trace_id attribute
	reqCtx := context.WithValue(ctx, TraceIDKey, traceID)
	reqLogger := logger.With(slog.String("trace_id", traceID))

	// 3. Pass enriched context to the metrics collector
	data, err := metrics.Collect(reqCtx, r.Config.Timeout)
	if err != nil {
		errMsg := "failed to collect metrics: " + err.Error()
		reqLogger.ErrorContext(reqCtx, errMsg, slog.String("error", err.Error()))
		sendToElastic("ERROR", errMsg, traceID)
		return
	}

	data.TraceID = traceID
	if data.Timestamp.IsZero() {
		data.Timestamp = time.Now().UTC()
	}

	if r.Storage != nil {
		r.Storage.Add(data)
	}

	metrics.RecordMetrics(data)

	hostname := data.System.System.Hostname
	if hostname == "" {
		hostname = "unknown"
	}

	msg := "metrics collected successfully"
	reqLogger.InfoContext(
		reqCtx,
		msg,
		slog.Time("timestamp", data.Timestamp),
		slog.String("hostname", hostname),
	)
	sendToElastic("INFO", msg, traceID)
}

func generateTraceID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "00000000000000000000000000000000"
	}
	return hex.EncodeToString(bytes)
}

func sendToElastic(level, message, traceID string) {
	esURL := os.Getenv("ELASTICSEARCH_URL")
	if esURL == "" {
		return
	}

	payload, err := json.Marshal(map[string]interface{}{
		"@timestamp": time.Now().UTC(),
		"level":      level,
		"message":    message,
		"trace_id":   traceID,
		"service":    "monitor-system",
	})
	if err != nil {
		return
	}

	go func() {
		client := &http.Client{Timeout: 2 * time.Second}
		_, _ = client.Post(esURL+"/monitor-system-logs/_doc", "application/json", bytes.NewBuffer(payload))
	}()
}
