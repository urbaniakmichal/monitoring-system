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

// HealthCheck godoc
// @Summary      Check agent status
// @Description  Returns the current agent running state (running/stopped) along with HATEOAS navigation links.
// @Tags         agent
// @Produce      json
// @Success      200  {object}  HealthResponse
// @Router       /api/v1/health [get]
func (rh *RestHandler) HealthCheck(res http.ResponseWriter, req *http.Request) {
	running := rh.as.IsRunning()

	statusStr := "stopped"
	relAction := "start"
	targetHref := ApiPathStart

	if running {
		statusStr = "running"
		relAction = "stop"
		targetHref = ApiPathStop
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
					Href:   ApiPathFile,
					Method: http.MethodGet,
				},
			},
		},
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(res).Encode(resp)
}

// StartAgent godoc
// @Summary      Start agent
// @Description  Starts background metrics collection loop.
// @Tags         agent
// @Produce      json
// @Success      200  {object}  AgentActionResponse
// @Failure      400  {object}  AgentActionResponse
// @Router       /api/v1/agent/start [post]
func (rh *RestHandler) StartAgent(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	if err := rh.as.Start(); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(res).Encode(AgentActionResponse{
			Message:   err.Error(),
			Timestamp: time.Now().UTC(),
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
					Href:   ApiPathStop,
					Method: http.MethodPost,
				},
				{
					Rel:    "metrics",
					Href:   ApiPathMetrics,
					Method: http.MethodGet,
				},
			},
		},
	}

	res.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(res).Encode(resp)
}

// StopAgent godoc
// @Summary      Stop agent
// @Description  Stops the active agent loop.
// @Tags         agent
// @Produce      json
// @Success      200  {object}  AgentActionResponse
// @Failure      400  {object}  AgentActionResponse
// @Router       /api/v1/agent/stop [post]
func (rh *RestHandler) StopAgent(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	if err := rh.as.Stop(); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(res).Encode(AgentActionResponse{
			Message:   err.Error(),
			Timestamp: time.Now().UTC(),
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
					Href:   ApiPathStart,
					Method: http.MethodPost,
				},
			},
		},
	}

	res.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(res).Encode(resp)
}

// Metrics godoc
// @Summary      Get collected metrics
// @Description  Returns an array of all collected metrics from cache.
// @Tags         agent
// @Produce      json
// @Success      200  {object}  MetricsResponse
// @Failure      500  {object}  AgentActionResponse
// @Router       /api/v1/agent/metrics [get]
func (rh *RestHandler) Metrics(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	metricsData, err := rh.as.Metrics()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(res).Encode(AgentActionResponse{
			Message:   err.Error(),
			Timestamp: time.Now().UTC(),
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
					Href:   ApiPathStart,
					Method: http.MethodPost,
				},
				{
					Rel:    "stop",
					Href:   ApiPathStop,
					Method: http.MethodPost,
				},
			},
		},
	}

	res.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(res).Encode(resp)
}

// GenerateFile godoc
// @Summary      Generate output file
// @Description  Triggers report file generation by the agent.
// @Tags         agent
// @Produce      multipart/form-data
// @Success      201  "Created"
// @Failure      400  {object}  map[string]string
// @Router       /api/v1/agent/file [get]
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
