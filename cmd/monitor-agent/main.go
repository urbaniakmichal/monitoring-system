package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"monitoring-system/internal/api"
	"monitoring-system/internal/config"
	"monitoring-system/internal/metrics"
	"monitoring-system/internal/runner"
	"monitoring-system/internal/storage"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("invalid config: %v", err)
	}

	store := storage.NewMemoryStorage(1000, logger)

	r := runner.Runner{
		Config:     *cfg,
		Collectors: metrics.NewCollectors(),
		Logger:     logger,
		Storage:    store,
	}

	ctrl := api.NewAgentController(r, ctx)
	handler := api.NewHandler(store, ctrl, logger)
	router := api.NewRouter(handler)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	ctrl.Start()

	go func() {
		logger.Info("starting HTTP server", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("HTTP server error", "error", err)
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down agent and HTTP server...")

	ctrl.Stop()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("HTTP server forced shutdown", "error", err)
	}

	logger.Info("agent exited gracefully")
}
