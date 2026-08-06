//go:build unit

package gpu

import (
	"errors"
	"testing"
)

func TestGpu_MockSuccess(t *testing.T) {
	mock := &MockGpu{
		DataToReturn: []GPUInformation{
			{
				Name:   "Name",
				Vendor: "Vendor",
			},
		},
	}

	_, err := mock.RetrieveGPUInfo()
	if err != nil {
		t.Fatalf("Absence of an error was expected: %v", err)
	}
}

func TestBattery_MockError(t *testing.T) {
	mock := &MockGpu{
		MockError: errors.New("unit test with mocks fail"),
	}

	_, err := mock.RetrieveGPUInfo()
	if err == nil {
		t.Fatal("Error was expected, got nil")
	}
}
