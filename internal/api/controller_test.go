package api_test

import (
	"context"
	"testing"
	"time"

	"monitoring-system/internal/api"
)

type DummyRunner struct{}

func (d *DummyRunner) Start(ctx context.Context) {
	<-ctx.Done()
}

func TestAgentController_StartAndStop(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := api.NewAgentController(&DummyRunner{}, ctx)

	if ctrl.IsRunning() {
		t.Errorf("expected running=false on startup")
	}

	if !ctrl.Start() {
		t.Errorf("expected Start() to return true")
	}

	time.Sleep(10 * time.Millisecond)

	if !ctrl.IsRunning() {
		t.Errorf("expected running=true after start")
	}

	if ctrl.Start() {
		t.Errorf("expected second Start() to return false")
	}

	if !ctrl.Stop() {
		t.Errorf("expected Stop() to return true")
	}

	if ctrl.IsRunning() {
		t.Errorf("expected running=false after stop")
	}
}

func TestAgentController_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	ctrl := api.NewAgentController(&DummyRunner{}, ctx)
	ctrl.Start()

	if !ctrl.IsRunning() {
		t.Fatalf("expected running=true")
	}

	cancel()
	time.Sleep(20 * time.Millisecond)

	if ctrl.IsRunning() {
		t.Errorf("expected running=false after parent context cancellation")
	}
}
