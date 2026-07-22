package metrics

import (
	mocks "monitor-agent/internal/testutil"
	"testing"
	"time"
)

func TestCollect_ShouldReturnSystemMetrics(t *testing.T) {

	providers := Collectors{

		Hostname: mocks.FakeHostname{
			Value: "test-hostname",
		},

		CPU: mocks.FakeCPU{
			Value: 50,
		},

		Memory: mocks.FakeMemory{
			Value: 50,
		},

		Disk: mocks.FakeDisk{
			Value: 50,
		},

		Processes: mocks.FakeProcesses{
			Value: []int32{1, 2, 3},
		},

		Uptime: mocks.FakeUptime{
			Value: 3600,
		},
	}

	result := Collect(providers, 2*time.Second)

	if result.Hostname == "" {
		t.Fatal("hostname should not be empty")
	}

	if result.CPU != 50 {
		t.Fatalf("invalid cpu: %d", result.CPU)
	}

	if result.Memory != 50 {
		t.Fatalf("invalid memory: %d", result.Memory)
	}

	if result.Disk != 50 {
		t.Fatalf("invalid disk: %d", result.Disk)
	}

	if result.Processes != 3 {
		t.Fatalf("invalid processes: %d", result.Processes)
	}

}
