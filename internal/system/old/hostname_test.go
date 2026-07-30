package system

import (
	"context"
	"errors"
	mocks "monitoring-system/internal/testutil"
	"testing"
	"time"
)

func TestHostname(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	reader := mocks.FakeHostname{Value: "test-hostname"}

	hostname, err := Hostname(ctx, reader)

	if err != nil {
		t.Fatalf("Hostname() returned error: %v", err)
	}

	if hostname == "" {
		t.Fatal("invalid hostname data")
	}
}

func TestHostname_Error(t *testing.T) {

	reader := mocks.FakeHostname{
		Err: errors.New("hostname failed"),
	}

	_, err := Hostname(context.Background(), reader)

	if err == nil {
		t.Fatal("expected error")
	}
}
