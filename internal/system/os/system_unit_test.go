//go:build unit

package system

import (
	"errors"
	"testing"
)

func TestSystemInformation_MockSuccess(t *testing.T) {
	mockSystemInformation := &MockSystemInformation{
   		DataToReturn: &SystemInformation{
			Hostname:     "Hostname",
			OS:           "OS",
			Architecture: "Architecture",
			Uptime:       "Uptime",
			TopProcesses: []ProcessInformation{
				{
					PID:        123,
					Name:       "Name",
					CPUPercent: 1.12,
					MemoryMB:   666,
				},
			},
		},
	}

	_, err := mockSystemInformation.RetrieveSystemInfo()
	if err != nil {
		t.Fatalf("Absence of an error was expected: %v", err)
	}
}

func TestSystemInformation_MockError(t *testing.T) {
	mockSystemInformation := &MockSystemInformation{
		MockError: errors.New("unit test with mocks fail"),
	}

	_, err := mockSystemInformation.RetrieveSystemInfo()
	if err == nil {
		t.Fatal("Error was expected, got nil")
	}
}