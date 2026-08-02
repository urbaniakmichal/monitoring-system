//go:build unit

package drivers

import (
	"errors"
	"testing"
)

func TestDrivers_MockSuccess(t *testing.T) {
	mockDriver := &MockDrivers{
		DataToReturn: []DriverInformation{
			{
				DeviceName:    "DeviceName",
				DriverVersion: "DriverVersion",
				Manufacturer:  "Manufacturer",
				DriverName:    "DriverName",
			},
		},
	}

	_, err := mockDriver.RetrieveInstalledDrivers()
	if err != nil {
		t.Fatalf("Absence of an error was expected: %v", err)
	}
}

func TestDrivers_MockError(t *testing.T) {
		mockDriver := &MockDrivers{
		MockError: errors.New("unit test with mocks fail"),
	}

	_, err := mockDriver.RetrieveInstalledDrivers()
	if err == nil {
		t.Fatal("Error was expected, got nil")
	}
}