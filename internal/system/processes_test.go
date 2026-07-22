package system

import (
	"context"
	mocks "monitor-agent/internal/testutil"
	"testing"
	"time"
)

func TestProcesses(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	reader := mocks.FakeProcesses{Value: []int32{1, 2, 3}}
	processes, err := Processes(ctx, reader)

	if err != nil {
		t.Fatalf("Processes() returned error: %v", err)
	}

	if processes < 0 {
		t.Fatalf("invalid processes value: %d", processes)
	}
}
