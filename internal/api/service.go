package api

import (
	"context"
	"errors"
	"log/slog"
	"sync"
)

type agentService struct {
	isRunning      bool
	mutex          sync.Mutex
	cancelFunction context.CancelFunc
	customLog      *slog.Logger
}

func NewAgentService(log *slog.Logger) *agentService {
	return &agentService{
		customLog: log,
	}
}

func (s *agentService) IsRunning() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.customLog.Info("agentService.isRunning value is", "isRunning", s.isRunning)
	return s.isRunning
}

func (s *agentService) Start() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.isRunning {
		return errors.New("Agent is already running")
	}

	s.isRunning = true
	s.customLog.Info("Running the service")
	return nil
}

func (s *agentService) Stop() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.isRunning {
		return errors.New("Agent is already stopped")
	}

	s.isRunning = false
	s.customLog.Info("Stopped the service")
	return nil
}

func (s *agentService) MakeFile() error {
	// todo implement later
	return nil
}
