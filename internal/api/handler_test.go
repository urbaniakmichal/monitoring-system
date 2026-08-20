//go:build api

package api

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"monitoring-system/internal/config"
	"monitoring-system/internal/memory_storage"
	"monitoring-system/internal/metrics"
	"monitoring-system/internal/runner"
)

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func newTestRunner(testLogger *slog.Logger) *runner.Runner {
	agentConfig := config.Config{
		Interval: 1 * time.Second,
	}

	metricsStorage := memory_storage.NewMemoryStorage(10, testLogger)

	return &runner.Runner{
		Config:  agentConfig,
		Storage: metricsStorage,
		Logger:  testLogger,
	}
}

func TestRequestHealthCheck(t *testing.T) {
	testCases := []struct {
		name                 string
		startServiceBefore   bool
		expectedHealthStatus string
		expectedFirstLinkRel string
	}{
		{
			name:                 "when agent is running",
			startServiceBefore:   true,
			expectedHealthStatus: "running",
			expectedFirstLinkRel: "stop",
		},
		{
			name:                 "when agent is stopped",
			startServiceBefore:   false,
			expectedHealthStatus: "stopped",
			expectedFirstLinkRel: "start",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testLogger := newTestLogger()
			agentRunner := newTestRunner(testLogger)
			agentService := NewAgentService(testLogger, agentRunner)
			restHandler := NewRestHandler(agentService)

			if testCase.startServiceBefore {
				// Uruchomienie przez handler, aby poprawnie zainicjalizować rh.startTime
				httpRequestStart := httptest.NewRequest(http.MethodPost, ApiPathStart, nil)
				recStart := httptest.NewRecorder()
				restHandler.StartAgent(recStart, httpRequestStart)
				time.Sleep(100 * time.Millisecond)
				defer agentService.Stop()
			}

			httpRequest := httptest.NewRequest(http.MethodGet, ApiPathHealth, nil)
			responseRecorder := httptest.NewRecorder()

			restHandler.HealthCheck(responseRecorder, httpRequest)
			httpResponse := responseRecorder.Result()
			defer httpResponse.Body.Close()

			if httpResponse.StatusCode != http.StatusOK {
				t.Fatalf("Expected status 200 OK, got %d", httpResponse.StatusCode)
			}

			var healthCheckResponse HealthResponse
			if err := json.NewDecoder(httpResponse.Body).Decode(&healthCheckResponse); err != nil {
				t.Fatalf("Failed to decode JSON: %v", err)
			}

			if healthCheckResponse.Status != testCase.expectedHealthStatus {
				t.Errorf("Expected status %q, got %q", testCase.expectedHealthStatus, healthCheckResponse.Status)
			}

			if len(healthCheckResponse.Links) < 2 {
				t.Fatalf("Expected at least 2 HATEOAS links, got %d", len(healthCheckResponse.Links))
			}

			if healthCheckResponse.Links[0].Rel != testCase.expectedFirstLinkRel {
				t.Errorf("Expected first link rel %q, got %q", testCase.expectedFirstLinkRel, healthCheckResponse.Links[0].Rel)
			}
			if healthCheckResponse.Links[1].Rel != "file" {
				t.Errorf("Expected second link rel 'file', got %q", healthCheckResponse.Links[1].Rel)
			}
		})
	}
}

func TestRequestStart(t *testing.T) {
	testCases := []struct {
		name                    string
		isAlreadyRunning        bool
		expectedStatusCode      int
		expectedResponseMessage string
		expectedLinkRel         string
	}{
		{
			name:                    "success - start stopped agent",
			isAlreadyRunning:        false,
			expectedStatusCode:      http.StatusOK,
			expectedResponseMessage: "agent started successfully",
			expectedLinkRel:         "stop",
		},
		{
			name:                    "failure - agent already running",
			isAlreadyRunning:        true,
			expectedStatusCode:      http.StatusBadRequest,
			expectedResponseMessage: "agent is already running",
			expectedLinkRel:         "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testLogger := newTestLogger()
			agentRunner := newTestRunner(testLogger)
			agentService := NewAgentService(testLogger, agentRunner)
			restHandler := NewRestHandler(agentService)

			if testCase.isAlreadyRunning {
				httpRequestStart := httptest.NewRequest(http.MethodPost, ApiPathStart, nil)
				recStart := httptest.NewRecorder()
				restHandler.StartAgent(recStart, httpRequestStart)
				time.Sleep(100 * time.Millisecond)
				defer agentService.Stop()
			}

			httpRequest := httptest.NewRequest(http.MethodPost, ApiPathStart, nil)
			responseRecorder := httptest.NewRecorder()

			restHandler.StartAgent(responseRecorder, httpRequest)
			httpResponse := responseRecorder.Result()
			defer httpResponse.Body.Close()

			if httpResponse.StatusCode != testCase.expectedStatusCode {
				t.Fatalf("Expected status %d, got %d", testCase.expectedStatusCode, httpResponse.StatusCode)
			}

			var actionResponse AgentActionResponse
			if err := json.NewDecoder(httpResponse.Body).Decode(&actionResponse); err != nil {
				t.Fatalf("Failed to decode JSON: %v", err)
			}

			if actionResponse.Message != testCase.expectedResponseMessage {
				t.Errorf("Expected message %q, got %q", testCase.expectedResponseMessage, actionResponse.Message)
			}

			if testCase.expectedLinkRel != "" {
				if len(actionResponse.Links) < 1 {
					t.Fatalf("Expected at least 1 HATEOAS link, got %d", len(actionResponse.Links))
				}
				if actionResponse.Links[0].Rel != testCase.expectedLinkRel {
					t.Errorf("Expected link rel %q, got %q", testCase.expectedLinkRel, actionResponse.Links[0].Rel)
				}
			}
		})
	}
}

func TestRequestStop(t *testing.T) {
	testCases := []struct {
		name                    string
		isAlreadyStopped        bool
		expectedStatusCode      int
		expectedResponseMessage string
		expectedLinkRel         string
	}{
		{
			name:                    "success - stop running agent",
			isAlreadyStopped:        false,
			expectedStatusCode:      http.StatusOK,
			expectedResponseMessage: "agent stopped successfully",
			expectedLinkRel:         "start",
		},
		{
			name:                    "failure - agent already stopped",
			isAlreadyStopped:        true,
			expectedStatusCode:      http.StatusBadRequest,
			expectedResponseMessage: "agent is already stopped",
			expectedLinkRel:         "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testLogger := newTestLogger()
			agentRunner := newTestRunner(testLogger)
			agentService := NewAgentService(testLogger, agentRunner)
			restHandler := NewRestHandler(agentService)

			if !testCase.isAlreadyStopped {
				httpRequestStart := httptest.NewRequest(http.MethodPost, ApiPathStart, nil)
				recStart := httptest.NewRecorder()
				restHandler.StartAgent(recStart, httpRequestStart)
				time.Sleep(100 * time.Millisecond)
				defer agentService.Stop()
			}

			httpRequest := httptest.NewRequest(http.MethodPost, ApiPathStop, nil)
			responseRecorder := httptest.NewRecorder()

			restHandler.StopAgent(responseRecorder, httpRequest)
			httpResponse := responseRecorder.Result()
			defer httpResponse.Body.Close()

			if httpResponse.StatusCode != testCase.expectedStatusCode {
				t.Fatalf("Expected status %d, got %d", testCase.expectedStatusCode, httpResponse.StatusCode)
			}

			var actionResponse AgentActionResponse
			if err := json.NewDecoder(httpResponse.Body).Decode(&actionResponse); err != nil {
				t.Fatalf("Failed to decode JSON: %v", err)
			}

			if actionResponse.Message != testCase.expectedResponseMessage {
				t.Errorf("Expected message %q, got %q", testCase.expectedResponseMessage, actionResponse.Message)
			}

			if testCase.expectedLinkRel != "" {
				if len(actionResponse.Links) < 1 {
					t.Fatalf("Expected at least 1 HATEOAS link, got %d", len(actionResponse.Links))
				}
				if actionResponse.Links[0].Rel != testCase.expectedLinkRel {
					t.Errorf("Expected link rel %q, got %q", testCase.expectedLinkRel, actionResponse.Links[0].Rel)
				}
			}
		})
	}
}

func TestRequestMetrics(t *testing.T) {
	testCases := []struct {
		name               string
		startServiceBefore bool
		isStorageNil       bool
		expectedStatusCode int
	}{
		{
			name:               "success - fetch metrics when agent is running",
			startServiceBefore: true,
			isStorageNil:       false,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "failure - storage not initialized when agent is stopped",
			startServiceBefore: false,
			isStorageNil:       true,
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testLogger := newTestLogger()
			agentRunner := newTestRunner(testLogger)
			if testCase.isStorageNil {
				agentRunner.Storage = nil
			}
			agentService := NewAgentService(testLogger, agentRunner)
			restHandler := NewRestHandler(agentService)

			if testCase.startServiceBefore {
				httpRequestStart := httptest.NewRequest(http.MethodPost, ApiPathStart, nil)
				recStart := httptest.NewRecorder()
				restHandler.StartAgent(recStart, httpRequestStart)
				time.Sleep(100 * time.Millisecond)
				defer agentService.Stop()
			}

			httpRequest := httptest.NewRequest(http.MethodGet, ApiPathMetrics, nil)
			responseRecorder := httptest.NewRecorder()

			restHandler.Metrics(responseRecorder, httpRequest)
			httpResponse := responseRecorder.Result()
			defer httpResponse.Body.Close()

			if httpResponse.StatusCode != testCase.expectedStatusCode {
				t.Fatalf("Expected status %d, got %d", testCase.expectedStatusCode, httpResponse.StatusCode)
			}

			if httpResponse.StatusCode == http.StatusOK {
				var metricsResponseBody MetricsResponse
				if err := json.NewDecoder(httpResponse.Body).Decode(&metricsResponseBody); err != nil {
					t.Fatalf("Failed to decode JSON: %v", err)
				}

				if metricsResponseBody.Data == nil {
					metricsResponseBody.Data = make([]metrics.Metrics, 0)
				}

				if len(metricsResponseBody.Links) == 0 {
					t.Errorf("Expected HATEOAS links in MetricsResponse")
				}
			}
		})
	}
}
