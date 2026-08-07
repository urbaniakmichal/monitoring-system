package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"monitoring-system/internal/api"
	"monitoring-system/internal/config"
	"monitoring-system/internal/logger"
	"monitoring-system/internal/memory_storage"
	"monitoring-system/internal/metrics"
	"monitoring-system/internal/runner"
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

	// Continuous background collection & HTTP server mode
	signalContext, signalCancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer signalCancel()

	memoryStorageInstance := memory_storage.NewMemoryStorage(1000, loggerInstance)

	agentRunner := runner.Runner{
		Config:  *cfg,
		Logger:  loggerInstance,
		Storage: memoryStorageInstance,
	}

	controllerInstance := api.NewAgentController(agentRunner, signalContext)
	httpServerInstance := buildHTTPServer(cfg, memoryStorageInstance, controllerInstance, loggerInstance)

	controllerInstance.Start()
	startHTTPServer(httpServerInstance, loggerInstance)

	loggerInstance.Info("starting monitoring agent application", "version", Version)

	// Wait for an interruption signal (blocks the main thread until shutdown)
	<-signalContext.Done()

	shutdownApplication(httpServerInstance, controllerInstance, loggerInstance)
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

func buildHTTPServer(cfg *config.Config, store *memory_storage.MemoryStorage, ctrl *api.AgentController, loggerInstance *slog.Logger) *http.Server {
	apiHandler := api.NewHandler(store, ctrl, loggerInstance)
	router := api.NewRouter(apiHandler)

	httpHandler := http.Handler(router)
	httpHandler = api.Logger(httpHandler)
	httpHandler = api.TraceID(httpHandler)
	httpHandler = api.Recovery(httpHandler)

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      httpHandler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}
}

func startHTTPServer(serverInstance *http.Server, loggerInstance *slog.Logger) {
	go func() {
		loggerInstance.Info("starting HTTP server", slog.String("addr", serverInstance.Addr))
		if serverError := serverInstance.ListenAndServe(); serverError != nil && !errors.Is(serverError, http.ErrServerClosed) {
			loggerInstance.Error("HTTP server error", slog.String("error_details", serverError.Error()))
		}
	}()
}

func shutdownApplication(serverInstance *http.Server, ctrl *api.AgentController, loggerInstance *slog.Logger) {
	loggerInstance.Info("shutting down agent and HTTP server...")

	ctrl.Stop()

	shutdownContext, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if shutdownError := serverInstance.Shutdown(shutdownContext); shutdownError != nil {
		loggerInstance.Error("HTTP server forced shutdown", slog.String("error_details", shutdownError.Error()))
	}

	loggerInstance.Info("agent exited gracefully")
}