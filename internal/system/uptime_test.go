package system

import (
	"context"
	"errors"
	mocks "monitoring-system/internal/testutil"
	"testing"
	"time"
)

func TestUptime(t *testing.T) {

	reader := mocks.FakeUptime{
		Value: 123,
	}

	uptime, err := Uptime(context.Background(), reader)

	if err != nil {
		t.Fatal(err)
	}

	expected := (123 * time.Second).String()

	if uptime != expected {
		t.Fatalf("expected %q, got %q", expected, uptime)
	}
}

func TestUptime_Error(t *testing.T) {

	reader := mocks.FakeUptime{
		Err: errors.New("uptime failed"),
	}

	_, err := Uptime(context.Background(), reader)

	if err == nil {
		t.Fatal("expected error")
	}
}
