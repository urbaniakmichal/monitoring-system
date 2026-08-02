//go:build unit

package operating_system

import (
	"errors"
	"testing"
)

func TestOperatingSystem_MockSuccess(t *testing.T) {
mockOperatingSystem := &MockOperatingSystem{
    DataToReturn: &OperatingSystemInformation{
        Caption:        "Caption",
        Version:        "Version",
        BuildNumber:    "BuildNumber",
        Manufacturer:   "Manufacturer",
        OSArchitecture: "OSArchitecture",
        InstallDate:    "InstallDate",
        LastBootTime:   "LastBootUpTime",
    },
}

	_, err := mockOperatingSystem.RetrieveOperatingSystemInformation()
	if err != nil {
		t.Fatalf("Absence of an error was expected: %v", err)
	}
}

func TestOperatingSystem_MockError(t *testing.T) {
	mockOperatingSystem := &MockOperatingSystem{
		MockError: errors.New("unit test with mocks fail"),
	}

	_, err := mockOperatingSystem.RetrieveOperatingSystemInformation()
	if err == nil {
		t.Fatal("Error was expected, got nil")
	}
}