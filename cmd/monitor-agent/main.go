package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"monitoring-system/internal/config"
	"monitoring-system/internal/logger"
	"monitoring-system/internal/metrics"
)

var Version = "dev"

func main() {
	// =========================================================================
	// CLI CHEAT SHEET & COMMAND EXAMPLES:
	// -------------------------------------------------------------------------
	// 1. Run as a continuous HTTP server (daemon mode):
	//    go run .\cmd\monitor-agent\main.go
	//    go run .\cmd\monitor-agent\main.go -config configs/config.yaml
	//
	// 2. Run in single-collection mode (CLI / One-shot) and print to console:
	//    go run .\cmd\monitor-agent\main.go -once -print-metrics=true
	//
	// 3. Run in single-collection mode and save output directly to a file:
	//    go run .\cmd\monitor-agent\main.go -once -output metrics.txt
	// =========================================================================

	// Define command-line flags
	printMetrics := flag.Bool("print-metrics", false, "Print collected metrics JSON to console in one-shot mode")
	onceFlag := flag.Bool("once", false, "Collect metrics once (on-demand) and exit")
	outputFile := flag.String("output", "", "Path to file where metrics JSON should be saved (used with -once)")
	configPath := flag.String("config", "configs/config.yaml", "Path to configuration file")
	flag.Parse()

	// Initialize logger and configuration
	loggerInstance := initLogger()
	cfg := loadConfig(*configPath)

	// Handle one-shot collection mode (on-demand / CLI)
	if *onceFlag {
		loggerInstance.Info("running single metrics collection (one-shot)...", "version", Version)

		collectionContext, collectionCancel := context.WithTimeout(context.Background(), cfg.Timeout)
		defer collectionCancel()

		data, err := metrics.Collect(collectionContext, cfg.Timeout)
		if err != nil {
			loggerInstance.Error("failed to collect single metrics", slog.String("error", err.Error()))
			os.Exit(1)
		}
		data.Timestamp = time.Now().UTC()

		// Format metrics to indented JSON
		jsonData, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			loggerInstance.Error("failed to marshal metrics to JSON", slog.String("error", err.Error()))
			os.Exit(1)
		}

		// If an output file path is specified, write JSON to the file
		if *outputFile != "" {
			err := os.WriteFile(*outputFile, jsonData, 0644)
			if err != nil {
				loggerInstance.Error("failed to write metrics to file", slog.String("file", *outputFile), slog.String("error", err.Error()))
				os.Exit(1)
			}
			loggerInstance.Info("successfully saved metrics to file", slog.String("file", *outputFile))
		}

		// Print to console if explicitly requested or if no output file was provided
		if *printMetrics || *outputFile == "" {
			fmt.Println(string(jsonData))
		}

		return
	}
}

func initLogger() *slog.Logger {
	logConfig := logger.Config{
		Level:  "debug",
		Format: "json",
	}

	log := logger.NewLogger(logConfig, os.Stdout)
	slog.SetDefault(log)

	return log
}

func loadConfig(path string) *config.Config {
	cfg, err := config.Load(path)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("invalid config: %v", err)
	}

	return cfg
}
