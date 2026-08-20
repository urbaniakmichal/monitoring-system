package runner

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
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
		reqLogger.ErrorContext(reqCtx, "failed to collect metrics", slog.String("error", err.Error()))
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

	reqLogger.InfoContext(
		reqCtx,
		"metrics collected successfully",
		slog.Time("timestamp", data.Timestamp),
		slog.String("hostname", hostname),
	)
}

func generateTraceID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "00000000000000000000000000000000"
	}
	return hex.EncodeToString(bytes)
}
