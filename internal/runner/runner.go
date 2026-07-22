package runner

import (
	"context"
	"log/slog"
	"monitor-agent/internal/config"
	"monitor-agent/internal/metrics"
	"time"
)

type Runner struct {
	Config     config.Config
	Collectors metrics.Collectors
	Logger     *slog.Logger
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
			data := metrics.Collect(r.Collectors, r.Config.Timeout)

			logger.Info(
				"metrics collected",
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
