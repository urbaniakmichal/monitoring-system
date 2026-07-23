package api

import "net/http"

func NewRouter(h *Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/metrics", h.GetMetrics)
	mux.HandleFunc("GET /api/v1/metrics/latest", h.GetLatestMetric)
	mux.HandleFunc("GET /api/v1/agent/status", h.GetStatus)
	mux.HandleFunc("POST /api/v1/agent/start", h.StartAgent)
	mux.HandleFunc("POST /api/v1/agent/stop", h.StopAgent)

	return mux
}
