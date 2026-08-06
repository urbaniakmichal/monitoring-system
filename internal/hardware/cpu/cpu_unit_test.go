//go:build unit

package cpu

import (
	"errors"
	"testing"
)


func TestCpuInformation_MockSuccess(t *testing.T) {
	mock := &MockCpuInformation{
   		DataToReturn: & CPUInformation{
			ModelName:    "ModelName",
			Cores:        123,
			UsagePercent: 1.23,
		},
	}

	_, err := mock.RetrieveCPUInfo()
	if err != nil {
		t.Fatalf("Absence of an error was expected: %v", err)
	}
}

func TestCpuInformation_MockError(t *testing.T) {
	mock := &MockCpuInformation{
		MockError: errors.New("unit test with mocks fail"),
	}

	_, err := mock.RetrieveCPUInfo()
	if err == nil {
		t.Fatal("Error was expected, got nil")
	}
}