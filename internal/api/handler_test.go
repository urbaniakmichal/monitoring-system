package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"monitoring-system/internal/api"
	"monitoring-system/internal/storage"
)

func TestHandler_StartAndStatusFlow(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	store := storage.NewMemoryStorage(10, nil)
	ctrl := api.NewAgentController(&DummyRunner{}, ctx)
	handler := api.NewHandler(store, ctrl, nil)
	router := api.NewRouter(handler)

	// 1. POST /api/v1/agent/start
	reqStart := httptest.NewRequest(http.MethodPost, "/api/v1/agent/start", nil)
	rrStart := httptest.NewRecorder()
	router.ServeHTTP(rrStart, reqStart)

	if rrStart.Code != http.StatusOK {
		t.Fatalf("POST /start returned status %d, expected %d", rrStart.Code, http.StatusOK)
	}

	// 2. GET /api/v1/agent/status – check if agent is STILL running after HTTP request completes
	reqStatus := httptest.NewRequest(http.MethodGet, "/api/v1/agent/status", nil)
	rrStatus := httptest.NewRecorder()
	router.ServeHTTP(rrStatus, reqStatus)

	if rrStatus.Code != http.StatusOK {
		t.Fatalf("GET /status returned status %d, expected %d", rrStatus.Code, http.StatusOK)
	}

	var resp struct {
		Data struct {
			Running bool `json:"running"`
		} `json:"data"`
	}

	if err := json.Unmarshal(rrStatus.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal JSON response: %v", err)
	}

	if !resp.Data.Running {
		t.Errorf("expected data.running = true after POST /start")
	}
}
