//go:build unit

package updates

import (
	"errors"
	"testing"
)

func TestUpdates_MockSuccess(t *testing.T) {
	mockApp := &MockUpdates{
		DataToReturn: []SystemUpdateInformation{
			{
				HotFixID:    "HotFixID",
				Description: "Description",
				InstalledOn: "InstalledOn",
				InstalledBy: "InstalledBy",
			},
		},
	}

	_, err := mockApp.RetrieveSystemUpdates()
	if err != nil {
		t.Fatalf("Absence of an error was expected: %v", err)
	}
}

func TestUpdates_MockError(t *testing.T) {
	mockApp := &MockUpdates{
		MockError: errors.New("unit test with mocks fail"),
	}

	_, err := mockApp.RetrieveSystemUpdates()
	if err == nil {
		t.Fatal("Error was expected, got nil")
	}
}
