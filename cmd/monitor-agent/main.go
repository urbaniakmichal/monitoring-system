package main

import (
	"errors"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"

	"monitoring-system/internal/api"
	"monitoring-system/internal/cli"
	"monitoring-system/internal/config"
	"monitoring-system/internal/logger"
	"monitoring-system/internal/memory_storage"
	"monitoring-system/internal/runner"
)

var Version = "dev"

func main() {

	// Initialize logger and configuration
	loggerInstance := initLogger()
	configPath := flag.String("config", "configs/config.yaml", "Path to configuration file")
	cfg := loadConfig(*configPath)
	run := &runner.Runner{
		Config:  *cfg,
		Logger:  loggerInstance,
		Storage: memory_storage.NewMemoryStorage(12, loggerInstance),
	}

	// CLI and flags
	printMetrics := flag.Bool("print-metrics", false, "Print collected metrics JSON to console in one-shot mode")
	onceFlag := flag.Bool("once", false, "Collect metrics once (on-demand) and exit")
	outputFile := flag.String("output", "", "Path to file where metrics JSON should be saved (used with -once)")
	flag.Parse()

	if *onceFlag {
		if err := cli.RunOneShot(loggerInstance, Version, cfg, *printMetrics, *outputFile); err != nil {
			loggerInstance.Error("cli execution failed", slog.String("error", err.Error()))
			os.Exit(1)
		}
		return
	}

	// Initialize server
	server := api.NewServer(loggerInstance, run, ":8080")
	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Server failed: %v", err)
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
