package storage

import (
	"log/slog"
	"monitoring-system/internal/metrics"
	"testing"
	"time"
)

func TestMemoryStorage_AddAndCapacity(t *testing.T) {
	s := NewMemoryStorage(2, slog.Default())

	s.Add(metrics.Metrics{CPU: 10})
	s.Add(metrics.Metrics{CPU: 20})
	s.Add(metrics.Metrics{CPU: 30})

	latest, ok := s.Latest()
	if !ok || latest.CPU != 30 {
		t.Fatalf("expected latest CPU 30, got %v", latest.CPU)
	}

	all := s.Query("cpu", "asc")
	if len(all) != 2 {
		t.Fatalf("expected capacity 2, got %d", len(all))
	}
	if all[0].CPU != 20 || all[1].CPU != 30 {
		t.Fatalf("unexpected items in buffer: %v", all)
	}
}

func TestMemoryStorage_QuerySorting(t *testing.T) {
	s := NewMemoryStorage(10, slog.Default())

	now := time.Now()
	s.Add(metrics.Metrics{Hostname: "B-server", CPU: 80, Timestamp: now.Add(1 * time.Second)})
	s.Add(metrics.Metrics{Hostname: "A-server", CPU: 20, Timestamp: now})

	sortedCPU := s.Query("cpu", "desc")
	if sortedCPU[0].CPU != 80 || sortedCPU[1].CPU != 20 {
		t.Fatalf("CPU desc sorting failed: %v", sortedCPU)
	}

	sortedHost := s.Query("hostname", "asc")
	if sortedHost[0].Hostname != "A-server" || sortedHost[1].Hostname != "B-server" {
		t.Fatalf("Hostname asc sorting failed: %v", sortedHost)
	}
}
