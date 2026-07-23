package api

import (
	"encoding/json"
	"log/slog"
	"monitoring-system/internal/storage"
	"net/http"
)

type Handler struct {
	storage    *storage.MemoryStorage
	controller *AgentController
	logger     *slog.Logger
}

func NewHandler(store *storage.MemoryStorage, ctrl *AgentController, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return &Handler{
		storage:    store,
		controller: ctrl,
		logger:     logger,
	}
}

func (h *Handler) buildLinks() map[string]Link {
	links := map[string]Link{
		"self":   {Href: "/api/v1/metrics", Method: "GET"},
		"latest": {Href: "/api/v1/metrics/latest", Method: "GET"},
		"status": {Href: "/api/v1/agent/status", Method: "GET"},
	}

	if h.controller.IsRunning() {
		links["stop_agent"] = Link{Href: "/api/v1/agent/stop", Method: "POST"}
	} else {
		links["start_agent"] = Link{Href: "/api/v1/agent/start", Method: "POST"}
	}

	return links
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("failed to encode JSON response", "error", err)
	}
}

// GET /api/v1/metrics?sort=cpu&order=desc
func (h *Handler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	sortBy := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	data := h.storage.Query(sortBy, order)

	resp := APIResponse{
		Data:  data,
		Links: h.buildLinks(),
	}
	h.writeJSON(w, http.StatusOK, resp)
}

// GET /api/v1/metrics/latest
func (h *Handler) GetLatestMetric(w http.ResponseWriter, r *http.Request) {
	latest, found := h.storage.Latest()
	if !found {
		h.writeJSON(w, http.StatusNotFound, ErrorResponse{
			Error: "no metrics collected yet",
			Links: h.buildLinks(),
		})
		return
	}

	resp := APIResponse{
		Data:  latest,
		Links: h.buildLinks(),
	}
	h.writeJSON(w, http.StatusOK, resp)
}

// GET /api/v1/agent/status
func (h *Handler) GetStatus(w http.ResponseWriter, r *http.Request) {
	status := map[string]any{
		"running": h.controller.IsRunning(),
	}

	resp := APIResponse{
		Data:  status,
		Links: h.buildLinks(),
	}
	h.writeJSON(w, http.StatusOK, resp)
}

// POST /api/v1/agent/start
func (h *Handler) StartAgent(w http.ResponseWriter, r *http.Request) {
	started := h.controller.Start()
	if !started {
		h.writeJSON(w, http.StatusConflict, ErrorResponse{
			Error: "agent is already running",
			Links: h.buildLinks(),
		})
		return
	}

	resp := APIResponse{
		Data:  map[string]string{"message": "agent started successfully"},
		Links: h.buildLinks(),
	}
	h.writeJSON(w, http.StatusOK, resp)
}

// POST /api/v1/agent/stop
func (h *Handler) StopAgent(w http.ResponseWriter, r *http.Request) {
	stopped := h.controller.Stop()
	if !stopped {
		h.writeJSON(w, http.StatusConflict, ErrorResponse{
			Error: "agent is not running",
			Links: h.buildLinks(),
		})
		return
	}

	resp := APIResponse{
		Data:  map[string]string{"message": "agent stopped successfully"},
		Links: h.buildLinks(),
	}
	h.writeJSON(w, http.StatusOK, resp)
}
