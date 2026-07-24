package api

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(h *Handler) http.Handler {
	mux := http.NewServeMux()

	// --- REST API Endpoints ---
	mux.HandleFunc("GET /api/v1/metrics", h.GetMetrics)
	mux.HandleFunc("GET /api/v1/metrics/latest", h.GetLatestMetric)
	mux.HandleFunc("GET /api/v1/agent/status", h.GetStatus)
	mux.HandleFunc("POST /api/v1/agent/start", h.StartAgent)
	mux.HandleFunc("POST /api/v1/agent/stop", h.StopAgent)

	// --- Infrastructure & Health Checks (Docker Compose / Kubernetes) ---
	mux.HandleFunc("GET /health", h.HealthCheck)

	// --- Telemetry & Metrics (Prometheus Exporter) ---
	mux.Handle("GET /metrics", promhttp.Handler())

	return mux
}
