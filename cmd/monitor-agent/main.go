package main

import (
	"context"
	"encoding/json"
	"errors"
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
	"monitoring-system/internal/metrics"
	"monitoring-system/internal/runner"
	"monitoring-system/internal/software"
	"monitoring-system/internal/system/old/storage"
)

var Version = "dev"

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	loggerInstance := initLogger()
	cfg := loadConfig("configs/config.yaml")

	store := storage.NewMemoryStorage(1000, loggerInstance)

	r := runner.Runner{
		Config:     *cfg,
		Collectors: metrics.NewCollectors(),
		Logger:     loggerInstance,
		Storage:    store,
	}

	ctrl := api.NewAgentController(r, ctx)
	server := buildHTTPServer(cfg, store, ctrl, loggerInstance)

	ctrl.Start()
	startServer(server, loggerInstance)

	<-ctx.Done()
	shutdown(server, ctrl, loggerInstance)

	/////////////////////////////////////////
	// new approach below - TODO change upper code in the near future

	slog.Info("----------------------------- NEW APPROACH ------------------------------------")
	slog.Info("Starting system monitoring application...")

	softwareData, collectionError := software.CollectAllSoftwareInformation()
	if collectionError != nil {
		slog.Error("Application failed to collect software data",
			slog.String("error_details", collectionError.Error()),
		)
		return
	}

	prettyJSON, marshallingError := json.MarshalIndent(softwareData, "", "  ")
	if marshallingError != nil {
		slog.Error("Failed to format collected data as JSON",
			slog.String("error_details", marshallingError.Error()),
		)
		return
	}

	fmt.Println(string(prettyJSON))
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

func buildHTTPServer(cfg *config.Config, store *storage.MemoryStorage, ctrl *api.AgentController, loggerInstance *slog.Logger) *http.Server {
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

func startServer(server *http.Server, loggerInstance *slog.Logger) {
	go func() {
		loggerInstance.Info("starting HTTP server", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			loggerInstance.Error("HTTP server error", "error", err)
		}
	}()
}

func shutdown(server *http.Server, ctrl *api.AgentController, loggerInstance *slog.Logger) {
	loggerInstance.Info("shutting down agent and HTTP server...")

	ctrl.Stop()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		loggerInstance.Error("HTTP server forced shutdown", "error", err)
	}

	loggerInstance.Info("agent exited gracefully")
}
