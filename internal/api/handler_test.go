//go:build api

package api

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestRequestHealthCheck(t *testing.T) {
	tests := []struct {
		name         string
		startService bool
		wantStatus   string
		wantRel      string
	}{
		{
			name:         "when agent is running",
			startService: true,
			wantStatus:   "running",
			wantRel:      "stop",
		},
		{
			name:         "when agent is stopped",
			startService: false,
			wantStatus:   "stopped",
			wantRel:      "start",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewAgentService(newTestLogger())
			if tt.startService {
				_ = svc.Start()
			}
			handler := NewRestHandler(svc)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
			w := httptest.NewRecorder()

			handler.HealthCheck(w, req)
			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != http.StatusOK {
				t.Fatalf("Expected status 200 OK, got %d", res.StatusCode)
			}

			var resp HealthResponse
			if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
				t.Fatalf("Failed to decode JSON: %v", err)
			}

			if resp.Status != tt.wantStatus {
				t.Errorf("Expected status %q, got %q", tt.wantStatus, resp.Status)
			}

			if len(resp.Links) < 2 {
				t.Fatalf("Expected at least 2 HATEOAS links, got %d", len(resp.Links))
			}

			if resp.Links[0].Rel != tt.wantRel {
				t.Errorf("Expected first link rel %q, got %q", tt.wantRel, resp.Links[0].Rel)
			}
			if resp.Links[1].Rel != "file" {
				t.Errorf("Expected second link rel 'file', got %q", resp.Links[1].Rel)
			}
		})
	}
}

func TestRequestStart(t *testing.T) {
	tests := []struct {
		name           string
		alreadyRunning bool
		wantStatusCode int
		wantMessage    string
		wantRel        string
	}{
		{
			name:           "success - start stopped agent",
			alreadyRunning: false,
			wantStatusCode: http.StatusOK,
			wantMessage:    "Agent started successfully",
			wantRel:        "stop",
		},
		{
			name:           "failure - agent already running",
			alreadyRunning: true,
			wantStatusCode: http.StatusBadRequest,
			wantMessage:    "Agent is already running",
			wantRel:        "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewAgentService(newTestLogger())
			if tt.alreadyRunning {
				_ = svc.Start()
			}
			handler := NewRestHandler(svc)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/agent/start", nil)
			w := httptest.NewRecorder()

			handler.StartAgent(w, req)
			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.wantStatusCode {
				t.Fatalf("Expected status %d, got %d", tt.wantStatusCode, res.StatusCode)
			}

			var resp AgentActionResponse
			if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
				t.Fatalf("Failed to decode JSON: %v", err)
			}

			if resp.Message != tt.wantMessage {
				t.Errorf("Expected message %q, got %q", tt.wantMessage, resp.Message)
			}

			if tt.wantRel != "" {
				if len(resp.Links) < 1 {
					t.Fatalf("Expected at least 1 HATEOAS link, got %d", len(resp.Links))
				}
				if resp.Links[0].Rel != tt.wantRel {
					t.Errorf("Expected link rel %q, got %q", tt.wantRel, resp.Links[0].Rel)
				}
			}
		})
	}
}

func TestRequestStop(t *testing.T) {
	tests := []struct {
		name           string
		alreadyStopped bool
		wantStatusCode int
		wantMessage    string
		wantRel        string
	}{
		{
			name:           "success - stop running agent",
			alreadyStopped: false,
			wantStatusCode: http.StatusOK,
			wantMessage:    "Agent stopped successfully",
			wantRel:        "start",
		},
		{
			name:           "failure - agent already stopped",
			alreadyStopped: true,
			wantStatusCode: http.StatusBadRequest,
			wantMessage:    "Agent is already stopped",
			wantRel:        "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewAgentService(newTestLogger())
			if !tt.alreadyStopped {
				_ = svc.Start()
			}
			handler := NewRestHandler(svc)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/agent/stop", nil)
			w := httptest.NewRecorder()

			handler.StopAgent(w, req)
			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.wantStatusCode {
				t.Fatalf("Expected status %d, got %d", tt.wantStatusCode, res.StatusCode)
			}

			var resp AgentActionResponse
			if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
				t.Fatalf("Failed to decode JSON: %v", err)
			}

			if resp.Message != tt.wantMessage {
				t.Errorf("Expected message %q, got %q", tt.wantMessage, resp.Message)
			}

			if tt.wantRel != "" {
				if len(resp.Links) < 1 {
					t.Fatalf("Expected at least 1 HATEOAS link, got %d", len(resp.Links))
				}
				if resp.Links[0].Rel != tt.wantRel {
					t.Errorf("Expected link rel %q, got %q", tt.wantRel, resp.Links[0].Rel)
				}
			}
		})
	}
}
