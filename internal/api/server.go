package api

import (
	"log/slog"
	"monitoring-system/internal/runner"
	"net/http"
	"time"
)

func NewServer(log *slog.Logger, run *runner.Runner, serverAddr string) *http.Server {
	agentSvc := NewAgentService(log, run)
	restHandler := NewRestHandler(agentSvc)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/health", restHandler.HealthCheck)
	mux.HandleFunc("POST /api/v1/agent/start", restHandler.StartAgent)
	mux.HandleFunc("POST /api/v1/agent/stop", restHandler.StopAgent)
	mux.HandleFunc("GET /api/v1/agent/file", restHandler.GenerateFile)
	mux.HandleFunc("GET /api/v1/agent/metrics", restHandler.Metrics)

	log.Info("server listening", slog.String("addr", serverAddr))

	return &http.Server{
		Addr:              serverAddr,
		Handler:           mux,
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       15 * time.Second,
	}
}
