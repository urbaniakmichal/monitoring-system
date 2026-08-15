package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"monitoring-system/internal/config"
	"monitoring-system/internal/metrics"
)

// =========================================================================
// CLI CHEAT SHEET & COMMAND EXAMPLES:
// -------------------------------------------------------------------------
//
//  1. Run as a continuous HTTP server (daemon mode):
//     go run .\cmd\monitor-agent\main.go
//     go run .\cmd\monitor-agent\main.go -config configs/config.yaml
//
//  2. Run in single-collection mode (CLI / One-shot) and print to console:
//     go run .\cmd\monitor-agent\main.go -once -print-metrics=true
//
//  3. Run in single-collection mode and save output directly to a file:
//     go run .\cmd\monitor-agent\main.go -once -output metrics.txt
//
// 4. Help go run .\cmd\monitor-agent\main.go -help
// =========================================================================
func RunOneShot(log *slog.Logger, version string, cfg *config.Config, printMetrics bool, outputFile string) error {
	log.Info("running single metrics collection (one-shot)...", "version", version)

	collectionContext, collectionCancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer collectionCancel()

	data, err := metrics.Collect(collectionContext, cfg.Timeout)
	if err != nil {
		return fmt.Errorf("collecting metrics failed: %w", err)
	}
	data.Timestamp = time.Now().UTC()

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling metrics to JSON failed: %w", err)
	}

	if outputFile != "" {
		if err := os.WriteFile(outputFile, jsonData, 0600); err != nil {
			return fmt.Errorf("writing metrics to file %s failed: %w", outputFile, err)
		}
		log.Info("successfully saved metrics to file", slog.String("file", outputFile))
	}

	if printMetrics || outputFile == "" {
		fmt.Println(string(jsonData))
	}

	return nil
}
