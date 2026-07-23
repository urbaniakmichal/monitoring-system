package runner

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"monitoring-system/internal/config"
	"monitoring-system/internal/metrics"
	"monitoring-system/internal/storage"
	"time"
)

type Runner struct {
	Config     config.Config
	Collectors metrics.Collectors
	Logger     *slog.Logger
	Storage    *storage.MemoryStorage
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

			data := metrics.Collect(r.Collectors, r.Config.Timeout)

			data.TraceID = traceID
			data.Timestamp = time.Now().UTC()

			if r.Storage != nil {
				r.Storage.Add(data)
			}

			logger.Info(
				"metrics collected",
				"trace_id", data.TraceID,
				"timestamp", data.Timestamp,
				"hostname", data.Hostname,
				"cpu", data.CPU,
				"memory", data.Memory,
				"disk", data.Disk,
				"uptime", data.Uptime,
				"processes", data.Processes,
			)
		}
	}
}
