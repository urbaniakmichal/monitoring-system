package memory_storage

import (
	"log/slog"
	"monitoring-system/internal/metrics"
	"sync"
)

type MemoryStorage struct {
	mtx      sync.RWMutex
	capacity int
	items    []metrics.Metrics
	logger   *slog.Logger
}

func NewMemoryStorage(capacity int, logger *slog.Logger) *MemoryStorage {
	return &MemoryStorage{
		capacity: capacity,
		items:    make([]metrics.Metrics, 0, capacity),
		logger:   logger,
	}
}

func (s *MemoryStorage) Add(item metrics.Metrics) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if len(s.items) >= s.capacity {
		// Remove the oldest item if capacity is reached (FIFO)
		s.items = s.items[1:]
	}
	s.items = append(s.items, item)

	if s.logger != nil {
		s.logger.Debug("metric added to storage", "total_items", len(s.items))
	}
}

func (s *MemoryStorage) GetAll() []metrics.Metrics {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	result := make([]metrics.Metrics, len(s.items))
	copy(result, s.items)
	return result
}
