package system

import (
	"context"
	"errors"
	mocks "monitor-agent/internal/testutil"
	"testing"
)

func TestMemory(t *testing.T) {

	reader := mocks.FakeMemory{
		Value: 50,
	}

	memory, err := Memory(context.Background(), reader)

	if err != nil {
		t.Fatalf("Memory() returned error: %v", err)
	}

	if memory < 0 || memory > 100 {
		t.Fatalf("invalid memory value: %d", memory)
	}
}

func TestMemory_Error(t *testing.T) {

	reader := mocks.FakeMemory{
		Err: errors.New("memory failed"),
	}

	_, err := Memory(context.Background(), reader)

	if err == nil {
		t.Fatal("expected error")
	}
}
