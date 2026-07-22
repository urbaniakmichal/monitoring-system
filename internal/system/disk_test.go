package system

import (
	"context"
	"errors"
	mocks "monitor-agent/internal/testutil"
	"testing"
	"time"
)

func TestDisk(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	reader := mocks.FakeDisk{Value: 10}
	disk, err := Disk(ctx, reader)

	if err != nil {
		t.Fatalf("Disk() returned error: %v", err)
	}

	if disk < 0 || disk > 100 {
		t.Fatal("invalid disk data")
	}
}

func TestDisk_Error(t *testing.T) {

	reader := mocks.FakeDisk{
		Err: errors.New("disk failed"),
	}

	_, err := Disk(context.Background(), reader)

	if err == nil {
		t.Fatal("expected error")
	}
}
