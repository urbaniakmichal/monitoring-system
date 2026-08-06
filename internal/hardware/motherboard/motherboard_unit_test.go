//go:build unit

package motherboard

import (
	"errors"
	"testing"
)

func TestMotherboard_MockSuccess(t *testing.T) {
	mock := &MockMotherboard{
   		DataToReturn: &MotherboardInformation{
			Manufacturer: "Manufacturer",
			Product:      "Product",
			Version:      "Version",
		},
	}

	_, err := mock.RetrieveMotherboardInfo()
	if err != nil {
		t.Fatalf("Absence of an error was expected: %v", err)
	}
}

func TestMotherboard_MockError(t *testing.T) {
	mock := &MockMotherboard{
		MockError: errors.New("unit test with mocks fail"),
	}

	_, err := mock.RetrieveMotherboardInfo()
	if err == nil {
		t.Fatal("Error was expected, got nil")
	}
}