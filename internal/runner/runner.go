package runner

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"monitoring-system/internal/config"
	"monitoring-system/internal/memory_storage"
	"monitoring-system/internal/metrics"
	"time"
)

type Runner struct {
	Config  config.Config
	Logger  *slog.Logger
	Storage *memory_storage.MemoryStorage
}

func generateTraceID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "00000000000000000000000000000000"
	}
	return hex.EncodeToString(bytes)
}

func (r Runner) Start(ctx context.Context) {
	logger := r.Logger
	if logger == nil {
		logger = slog.Default()
	}

	ticker := time.NewTicker(r.Config.Interval)
	defer ticker.Stop()

	logger.Info("starting metrics collector", "interval", r.Config.Interval)

	for {
		select {
		case <-ctx.Done():
			logger.Info("stopping metrics collector")
			return
		case <-ticker.C:
			traceID := generateTraceID()

			data, err := metrics.Collect(ctx, r.Config.Timeout)
			if err != nil {
				logger.Error("failed to collect metrics", "error", err)
				continue
			}

			data.TraceID = traceID
			data.Timestamp = time.Now().UTC()

			if r.Storage != nil {
				r.Storage.Add(data)
			}

			metrics.RecordMetrics(data)

			hostname := data.System.System.Hostname
			if hostname == "" {
				hostname = "unknown"
			}

			logger.Info(
				"metrics collected successfully",
				"trace_id", data.TraceID,
				"timestamp", data.Timestamp,
				"hostname", hostname,
			)
		}
	}
}
