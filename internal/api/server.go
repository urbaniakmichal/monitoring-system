package api

import (
	"fmt"
	"log/slog"
	"monitoring-system/internal/config"
	"monitoring-system/internal/runner"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

func NewServer(log *slog.Logger, run *runner.Runner, cfg config.ServerConfig) *http.Server {
	agentSvc := NewAgentService(log, run)
	restHandler := NewRestHandler(agentSvc)

	mux := http.NewServeMux()
	mux.HandleFunc("GET "+ApiPathHealth, restHandler.HealthCheck)
	mux.HandleFunc("POST "+ApiPathStart, restHandler.StartAgent)
	mux.HandleFunc("POST "+ApiPathStop, restHandler.StopAgent)
	mux.HandleFunc("GET "+ApiPathFile, restHandler.GenerateFile)
	mux.HandleFunc("GET "+ApiPathMetrics, restHandler.Metrics)

	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	serverAddr := fmt.Sprintf(":%d", cfg.Port)

	return &http.Server{
		Addr:              serverAddr,
		Handler:           mux,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
	}
}
