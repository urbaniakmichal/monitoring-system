//go:build unit

package startup

import (
	"errors"
	"testing"
)

func TestStartup_MockSuccess(t *testing.T) {
	mockApp := &MockStartup{
		DataToReturn: []StartupCommandInformation{
			{
				Name:     "Name",
				Command:  "Command",
				Location: "Location",
				User:     "User",
			},
		},
	}

	_, err := mockApp.RetrieveStartupCommands()
	if err != nil {
		t.Fatalf("Absence of an error was expected: %v", err)
	}
}

func TestStartup_MockError(t *testing.T) {
	mockApp := &MockStartup{
		MockError: errors.New("unit test with mocks fail"),
	}

	_, err := mockApp.RetrieveStartupCommands()
	if err == nil {
		t.Fatal("Error was expected, got nil")
	}
}