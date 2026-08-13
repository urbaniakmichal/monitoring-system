package api

import (
	"encoding/json"
	"net/http"
	"time"
)

type RestHandler struct {
	as *AgentService
}

func NewRestHandler(as *AgentService) *RestHandler {
	return &RestHandler{
		as: as,
	}
}

func (rh *RestHandler) HealthCheck(res http.ResponseWriter, req *http.Request) {
	running := rh.as.IsRunning()

	statusStr := "stopped"
	relAction := "start"
	targetHref := "/api/v1/agent/start"

	if running {
		statusStr = "running"
		relAction = "stop"
		targetHref = "/api/v1/agent/stop"
	}

	resp := HealthResponse{
		Status:    statusStr,
		Uptime:    "active",
		Timestamp: time.Now().UTC(),
		ResponseEnvelope: ResponseEnvelope{
			Links: []Link{
				{
					Rel:    relAction,
					Href:   targetHref,
					Method: http.MethodPost,
				},
				{
					Rel:    "file",
					Href:   "/api/v1/agent/file",
					Method: http.MethodGet,
				},
			},
		},
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(res).Encode(resp)
}

func (rh *RestHandler) StartAgent(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	if err := rh.as.Start(); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(res).Encode(AgentActionResponse{
			Message: err.Error(),
		})
		return
	}

	resp := AgentActionResponse{
		Message:   "agent started successfully",
		Timestamp: time.Now().UTC(),
		ResponseEnvelope: ResponseEnvelope{
			Links: []Link{
				{
					Rel:    "stop",
					Href:   "/api/v1/agent/stop",
					Method: http.MethodPost,
				},
				{
					Rel:    "metrics",
					Href:   "/api/v1/agent/metrics",
					Method: http.MethodGet,
				},
			},
		},
	}

	res.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(res).Encode(resp)
}

func (rh *RestHandler) StopAgent(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	if err := rh.as.Stop(); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(res).Encode(AgentActionResponse{
			Message: err.Error(),
		})
		return
	}

	resp := AgentActionResponse{
		Message:   "agent stopped successfully",
		Timestamp: time.Now().UTC(),
		ResponseEnvelope: ResponseEnvelope{
			Links: []Link{
				{
					Rel:    "start",
					Href:   "/api/v1/agent/start",
					Method: http.MethodPost,
				},
			},
		},
	}

	res.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(res).Encode(resp)
}

func (rh *RestHandler) Metrics(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	metricsData, err := rh.as.Metrics()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(res).Encode(AgentActionResponse{
			Message: err.Error(),
		})
		return
	}

	resp := MetricsResponse{
		Message:   "get all metrics successfully",
		Timestamp: time.Now().UTC(),
		Data:      metricsData,
		ResponseEnvelope: ResponseEnvelope{
			Links: []Link{
				{
					Rel:    "start",
					Href:   "/api/v1/agent/start",
					Method: http.MethodPost,
				},
				{
					Rel:    "stop",
					Href:   "/api/v1/agent/stop",
					Method: http.MethodPost,
				},
			},
		},
	}

	res.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(res).Encode(resp)
}

func (rh *RestHandler) GenerateFile(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "multipart/form-data")

	err := rh.as.MakeFile()
	if err != nil {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(res).Encode(map[string]string{"error": err.Error()})
		return
	}

	res.WriteHeader(http.StatusCreated)
}
