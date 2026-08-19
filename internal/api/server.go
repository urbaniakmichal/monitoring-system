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
		Handler:           enableCORS(mux),
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
	}
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
