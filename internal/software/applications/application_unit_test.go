//go:build unit

package applications

import (
	"errors"
	"testing"
)

func TestApplications_MockSuccess(t *testing.T) {
	mockApp := &MockApplications{
		AppsToReturn: []ApplicationInformation{
			{
				DisplayName:    "Mock display name",
				DisplayVersion: "Mock display version",
				Publisher:      "Mock Publisher",
				InstallDate:    "20/02/2000",
			},
		},
	}

	_, err := mockApp.RetrieveInstalledApplications()
	if err != nil {
		t.Fatalf("Absence of an error was expected: %v", err)
	}
}

func TestApplications_MockError(t *testing.T) {
	mockApp := &MockApplications{
		MockError: errors.New("unit test with mocks fail"),
	}

	_, err := mockApp.RetrieveInstalledApplications()
	if err == nil {
		t.Fatal("Error was expected, got nil")
	}
}