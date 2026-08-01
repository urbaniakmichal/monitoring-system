package runner

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"monitoring-system/internal/config"
	"monitoring-system/internal/metrics"
)

// Runner coordinates the metrics collection process.
type Runner struct {
	Config       config.Config
	Logger       *slog.Logger
	PrintMetrics bool
}

// generateTraceID creates a unique trace identifier for a collection cycle.
func generateTraceID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "00000000000000000000000000000000"
	}
	return hex.EncodeToString(bytes)
}

// CollectOnce gathers system metrics a single time on-demand.
func (r Runner) CollectOnce(ctx context.Context) (metrics.Metrics, error) {
	traceID := generateTraceID()

	data, err := metrics.Collect(ctx, r.Config.Timeout)
	if err != nil {
		return metrics.Metrics{}, err
	}

	data.TraceID = traceID
	data.Timestamp = time.Now().UTC()

	// Update metrics on Prometheus if applicable
	metrics.RecordMetrics(data)

	return data, nil
}

// Start runs the continuous metrics collection loop in the background.
func (r Runner) Start(ctx context.Context) {
	logger := r.Logger
	if logger == nil {
		logger = slog.Default()
	}

	ticker := time.NewTicker(r.Config.Interval)
	defer ticker.Stop()

	logger.Info("starting comprehensive metrics collector", "interval", r.Config.Interval, "print_metrics", r.PrintMetrics)

	for {
		select {
		case <-ctx.Done():
			logger.Info("stopping metrics collector")
			return
		case <-ticker.C:
			data, err := r.CollectOnce(ctx)
			if err != nil {
				logger.Error("failed to collect system metrics", "error", err)
				continue
			}

			hostname := data.System.System.Hostname
			if hostname == "" {
				hostname = "unknown"
			}

			logger.Info(
				"comprehensive metrics collected successfully",
				"trace_id", data.TraceID,
				"timestamp", data.Timestamp,
				"hostname", hostname,
			)

			// Print metrics JSON to the console if the flag is enabled e.g.  go run .\cmd\monitor-agent\main.go -once -print-metrics=true -output metrics.txt
			if r.PrintMetrics {
				jsonData, err := json.MarshalIndent(data, "", "  ")
				if err == nil {
					fmt.Println("\n=== [DEBUG] COLLECTED METRICS JSON START ===")
					fmt.Println(string(jsonData))
					fmt.Println("\n=== [DEBUG] COLLECTED METRICS JSON END ===")
				} else {
					logger.Error("Failed to marshal metrics for console print", "error", err)
				}
			}
		}
	}
}
