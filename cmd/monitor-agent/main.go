package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "monitoring-system/docs" // Register generated Swagger UI documentation

	"monitoring-system/internal/api"
	"monitoring-system/internal/cli"
	"monitoring-system/internal/config"
	"monitoring-system/internal/logger"
	"monitoring-system/internal/memory_storage"
	"monitoring-system/internal/runner"
)

var Version = "dev"

// @title           Monitor Agent API
// @version         1.0
// @description     REST API for system, hardware, and software monitoring.

// @host            localhost:8080
// @BasePath        /
func main() {
	// 1. Define command-line flags
	configPath := flag.String("config", "configs/config.yaml", "Path to configuration file")
	printMetrics := flag.Bool("print-metrics", false, "Print collected metrics JSON to console in one-shot mode")
	onceFlag := flag.Bool("once", false, "Collect metrics once (on-demand) and exit")
	outputFile := flag.String("output", "", "Path to file where metrics JSON should be saved (used with -once)")

	// 2. Parse flags after defining them
	flag.Parse()

	// 3. Load configuration
	cfg := loadConfig(*configPath)

	// 4. Initialize components using loaded configuration
	loggerInstance := initLogger(cfg.Logger)
	run := &runner.Runner{
		Config:  *cfg,
		Logger:  loggerInstance,
		Storage: memory_storage.NewMemoryStorage(cfg.Storage.RetentionHours, loggerInstance),
	}

	if *onceFlag {
		if err := cli.RunOneShot(loggerInstance, Version, cfg, *printMetrics, *outputFile); err != nil {
			loggerInstance.Error("cli execution failed", slog.String("error", err.Error()))
			os.Exit(1)
		}
		return
	}

	// 5. Start HTTP server and handle graceful shutdown
	server := api.NewServer(loggerInstance, run, cfg.Server)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	serverErrChan := make(chan error, 1)

	go func() {
		loggerInstance.Info("starting server", slog.String("port", fmt.Sprintf(":%d", cfg.Server.Port)))

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
		loggerInstance.Info("initiating graceful shutdown...")

		shutdownCtx, cancelCtx := context.WithTimeout(context.Background(), cfg.Timeout)
		defer cancelCtx()

		if err := server.Shutdown(shutdownCtx); err != nil {
			loggerInstance.Error("server shutdown error", slog.String("error", err.Error()))
		}
	}

	loggerInstance.Info("server shutdown complete")
}

func initLogger(cfg config.LoggerConfig) *slog.Logger {
	logConfig := logger.Config{
		Level:  cfg.Level,
		Format: cfg.Format,
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
