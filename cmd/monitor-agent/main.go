package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"monitoring-system/internal/api"
	"monitoring-system/internal/cli"
	"monitoring-system/internal/config"
	"monitoring-system/internal/logger"
	"monitoring-system/internal/memory_storage"
	"monitoring-system/internal/runner"
)

var Version = "dev"

func main() {

	// CLI and flags
	flag.Parse()
	configPath := flag.String("config", "configs/config.yaml", "Path to configuration file")
	printMetrics := flag.Bool("print-metrics", false, "Print collected metrics JSON to console in one-shot mode")
	onceFlag := flag.Bool("once", false, "Collect metrics once (on-demand) and exit")
	outputFile := flag.String("output", "", "Path to file where metrics JSON should be saved (used with -once)")

	// Initialize logger and configuration
	loggerInstance := initLogger()
	cfg := loadConfig(*configPath)
	run := &runner.Runner{
		Config:  *cfg,
		Logger:  loggerInstance,
		Storage: memory_storage.NewMemoryStorage(12, loggerInstance),
	}

	if *onceFlag {
		if err := cli.RunOneShot(loggerInstance, Version, cfg, *printMetrics, *outputFile); err != nil {
			loggerInstance.Error("cli execution failed", slog.String("error", err.Error()))
			os.Exit(1)
		}
		return
	}

	// Server
	server := api.NewServer(loggerInstance, run, ":8080")
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	serverErrChan := make(chan error, 1)

	go func() {
		loggerInstance.Info("starting server", slog.String("port", ":8080"))

		err := server.ListenAndServe()

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrChan <- err
		}
	}()

	select {
	case err := <-serverErrChan:
		loggerInstance.Error("server failed to start", slog.String("error", err.Error()))
		os.Exit(1)

	case <-ctx.Done():
		loggerInstance.Info("initiating graceful shutdown start... (e.g. because Ctrl+C / SIGTERM )")

		shutdownCtx, cancelCtx := context.WithTimeout(context.Background(), cfg.Timeout)
		defer cancelCtx()
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			loggerInstance.Error("server shutdown error", slog.String("error", err.Error()))
		}

		loggerInstance.Info("initiating graceful shutdown done")
	}

	loggerInstance.Info("server shutdown complete")
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
