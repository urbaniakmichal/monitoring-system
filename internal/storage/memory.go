package storage

import (
	"log/slog"
	"monitoring-system/internal/metrics"
	"slices"
	"strings"
	"sync"
)

type MemoryStorage struct {
	mu       sync.RWMutex
	metrics  []metrics.Metrics
	capacity int
	logger   *slog.Logger
}

func NewMemoryStorage(capacity int, logger *slog.Logger) *MemoryStorage {
	if logger == nil {
		logger = slog.Default()
	}

	return &MemoryStorage{
		metrics:  make([]metrics.Metrics, 0, capacity),
		capacity: capacity,
		logger:   logger,
	}
}

func (s *MemoryStorage) Add(m metrics.Metrics) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.metrics) >= s.capacity {
		s.logger.Debug("storage buffer full, evicting oldest metric",
			"capacity", s.capacity,
			"evicted_trace_id", s.metrics[0].TraceID,
			"evicted_timestamp", s.metrics[0].Timestamp,
		)
		s.metrics = s.metrics[1:]
	}

	s.metrics = append(s.metrics, m)

	s.logger.Debug("metric added to storage",
		"trace_id", m.TraceID,
		"hostname", m.Hostname,
		"cpu", m.CPU,
		"memory", m.Memory,
		"total_items", len(s.metrics),
	)
}

func (s *MemoryStorage) Latest() (metrics.Metrics, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.metrics) == 0 {
		s.logger.Debug("latest metric requested but storage is empty")
		return metrics.Metrics{}, false
	}

	return s.metrics[len(s.metrics)-1], true
}

func (s *MemoryStorage) Query(sortBy, order string) []metrics.Metrics {
	s.mu.RLock()
	defer s.mu.RUnlock()

	s.logger.Debug("querying metrics",
		"sort_by", sortBy,
		"order", order,
		"total_available", len(s.metrics),
	)

	if len(s.metrics) == 0 {
		return []metrics.Metrics{}
	}

	result := make([]metrics.Metrics, len(s.metrics))
	copy(result, s.metrics)

	sortBy = strings.ToLower(sortBy)
	order = strings.ToLower(order)

	slices.SortFunc(result, func(a, b metrics.Metrics) int {
		var cmp int
		switch sortBy {
		case "cpu":
			cmp = cmpInt(a.CPU, b.CPU)
		case "memory":
			cmp = cmpInt(a.Memory, b.Memory)
		case "hostname":
			cmp = strings.Compare(a.Hostname, b.Hostname)
		case "timestamp":
			fallthrough
		default:
			if a.Timestamp.Before(b.Timestamp) {
				cmp = -1
			} else if a.Timestamp.After(b.Timestamp) {
				cmp = 1
			} else {
				cmp = 0
			}
		}

		if order == "desc" {
			return -cmp
		}
		return cmp
	})

	return result
}

func cmpInt(a, b int) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}
