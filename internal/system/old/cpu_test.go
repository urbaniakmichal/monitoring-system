package system

import (
	"context"
	"errors"
	mocks "monitoring-system/internal/testutil"
	"testing"
)

func TestCPU(t *testing.T) {

	reader := mocks.FakeCPU{Value: 50}

	cpu, err := CPU(context.Background(), reader)

	if err != nil {
		t.Fatal(err)
	}

	if cpu != 50 {
		t.Fatal("invalid cpu data")
	}
}

func TestCPU_Error(t *testing.T) {

	reader := mocks.FakeCPU{
		Err: errors.New("cpu failed"),
	}

	_, err := CPU(context.Background(), reader)

	if err == nil {
		t.Fatal("expected error")
	}
}
