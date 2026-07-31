package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"monitoring-system/internal/config"
	"monitoring-system/internal/logger"
	"monitoring-system/internal/runner"
)

var Version = "dev"

func main() {
	// Define command-line flags
	printMetrics := flag.Bool("print-metrics", false, "Print collected metrics JSON to console in loop mode")
	onceFlag := flag.Bool("once", false, "Collect metrics once (on-demand) and exit")
	flag.Parse()

	// Context to handle graceful shutdown (Ctrl+C, SIGTERM)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Initialize logger and configuration
	loggerInstance := initLogger()
	cfg := loadConfig("configs/config.yaml")

	// Initialize the runner
	r := runner.Runner{
		Config:       *cfg,
		Logger:       loggerInstance,
		PrintMetrics: *printMetrics,
	}

	// Handle one-shot collection mode (on-demand, e.g., for future API use)
	if *onceFlag {
		loggerInstance.Info("running single metrics collection (one-shot)...", "version", Version)

		data, err := r.CollectOnce(ctx)
		if err != nil {
			loggerInstance.Error("failed to collect single metrics", "error", err)
			os.Exit(1)
		}

		// Output formatted JSON to stdout
		jsonData, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			loggerInstance.Error("failed to marshal metrics to JSON", "error", err)
			os.Exit(1)
		}

		fmt.Println(string(jsonData))
		return
	}

	// Continuous background collection mode
	loggerInstance.Info("starting monitoring agent application", "version", Version)
	go r.Start(ctx)

	// Wait for an interruption signal (blocks the main thread until shutdown)
	<-ctx.Done()

	loggerInstance.Info("shutting down monitoring agent gracefully...")
	loggerInstance.Info("agent exited")
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