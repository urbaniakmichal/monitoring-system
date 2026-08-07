package memory_storage

import (
	"log/slog"
	"sync"
)

type MemoryStorage struct {
	mu       sync.RWMutex
	capacity int
	items    []any
	logger   *slog.Logger
}

func NewMemoryStorage(capacity int, logger *slog.Logger) *MemoryStorage {
	return &MemoryStorage{
		capacity: capacity,
		items:    make([]any, 0, capacity),
		logger:   logger,
	}
}

func (s *MemoryStorage) Add(item any) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.items) >= s.capacity {
		// Remove the oldest item if capacity is reached (FIFO)
		s.items = s.items[1:]
	}
	s.items = append(s.items, item)

	if s.logger != nil {
		s.logger.Debug("metric added to storage", "total_items", len(s.items))
	}
}

func (s *MemoryStorage) GetAll() []any {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]any, len(s.items))
	copy(result, s.items)
	return result
}