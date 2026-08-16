package api

import (
	"context"
	"errors"
	"log/slog"
	"monitoring-system/internal/metrics"
	"monitoring-system/internal/runner"
	"sync"
)

type AgentService struct {
	isRunning      bool
	mutex          sync.RWMutex
	cancelFunction context.CancelFunc
	wg             sync.WaitGroup
	customLog      *slog.Logger
	runner         *runner.Runner
}

func NewAgentService(log *slog.Logger, r *runner.Runner) *AgentService {
	return &AgentService{
		customLog: log,
		runner:    r,
	}
}

func (s *AgentService) IsRunning() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	s.customLog.Info("AgentService.isRunning value is", "isRunning", s.isRunning)
	return s.isRunning
}

func (s *AgentService) Start() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.isRunning {
		return errors.New("agent is already running")
	}

	childCtx, cancel := context.WithCancel(context.Background())
	s.cancelFunction = cancel
	go s.runner.Start(childCtx)

	s.isRunning = true
	s.customLog.Info("Running the service")

	return nil
}

func (s *AgentService) Stop() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.isRunning {
		return errors.New("agent is already stopped")
	}

	if s.cancelFunction != nil {
		s.cancelFunction()
	}

	s.wg.Wait()

	s.cancelFunction = nil
	s.isRunning = false
	s.customLog.Info("Stopped the service")

	return nil
}

func (s *AgentService) Metrics() ([]metrics.Metrics, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if s.runner == nil || s.runner.Storage == nil {
		return nil, errors.New("storage is not initialized")
	}

	return s.runner.Storage.GetAll(), nil
}

func (s *AgentService) MakeFile() error {
	// todo implement later
	return nil
}
