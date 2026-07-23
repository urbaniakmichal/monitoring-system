package api

import (
	"context"
	"sync"
)

type Runner interface {
	Start(ctx context.Context)
}

type AgentController struct {
	mu      sync.Mutex
	runner  Runner
	cancel  context.CancelFunc
	running bool
	baseCtx context.Context
}

func NewAgentController(r Runner, baseCtx context.Context) *AgentController {
	if baseCtx == nil {
		baseCtx = context.Background()
	}
	return &AgentController{
		runner:  r,
		running: false,
		baseCtx: baseCtx,
	}
}

func (c *AgentController) Start() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running {
		return false
	}

	ctx, cancel := context.WithCancel(c.baseCtx)
	c.cancel = cancel
	c.running = true

	go func() {
		c.runner.Start(ctx)
		c.mu.Lock()
		c.running = false
		c.mu.Unlock()
	}()

	return true
}

func (c *AgentController) Stop() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.running || c.cancel == nil {
		return false
	}

	c.cancel()
	c.running = false
	return true
}

func (c *AgentController) IsRunning() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.running
}
