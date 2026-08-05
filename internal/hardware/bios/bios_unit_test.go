//go:build unit

package bios

import (
	"errors"
	"testing"
)

func TestBios_MockSuccess(t *testing.T) {
	mock := &MockBios{
   		DataToReturn: &BiosInformation{
			Manufacturer: "Manufacturer",
			Version:      "Version",
			ReleaseDate:  "ReleaseDate",
			SerialNumber: "SerialNumber",
		},
	}

	_, err := mock.RetrieveBiosInformation()
	if err != nil {
		t.Fatalf("Absence of an error was expected: %v", err)
	}
}

func TestBios_MockError(t *testing.T) {
	mock := &MockBios{
		MockError: errors.New("unit test with mocks fail"),
	}

	_, err := mock.RetrieveBiosInformation()
	if err == nil {
		t.Fatal("Error was expected, got nil")
	}
}