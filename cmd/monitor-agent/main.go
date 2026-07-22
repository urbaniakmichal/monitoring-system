package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"monitoring-system/internal/config"
	"monitoring-system/internal/metrics"
	"monitoring-system/internal/runner"
)

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("invalid config: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	r := runner.Runner{
		Config:     *cfg,
		Collectors: metrics.NewCollectors(),
		Logger:     slog.Default(),
	}

	r.Start(ctx)
}
