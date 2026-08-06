//go:build unit

package memory

import (
	"errors"
	"testing"
)

func TestMemory_MockSuccess(t *testing.T) {
	mock := &MockMemory{
   		DataToReturn: &MemoryInformation{
				TotalMB:     12,
			AvailableMB: 13,
			UsedPercent: 1.23,
		},
	}

	_, err := mock.RetrieveMemoryInfo()
	if err != nil {
		t.Fatalf("Absence of an error was expected: %v", err)
	}
}

func TestMemory_MockError(t *testing.T) {
	mock := &MockMemory{
		MockError: errors.New("unit test with mocks fail"),
	}

	_, err := mock.RetrieveMemoryInfo()
	if err == nil {
		t.Fatal("Error was expected, got nil")
	}
}