//go:build unit

package services

import (
	"errors"
	"testing"
)

func TestServices_MockSuccess(t *testing.T) {
	mockApp := &MockServices{
		DataToReturn: []ServiceInformation{
			{
				Name:        "Name",
				DisplayName: "DisplayName",
				State:       "State",
				StartMode:   "StartMode",
				StartName:   "StartName",
			},
		},
	}

	_, err := mockApp.RetrieveSystemServices()
	if err != nil {
		t.Fatalf("Absence of an error was expected: %v", err)
	}
}

func TestServices_MockError(t *testing.T) {
	mockApp := &MockServices{
		MockError: errors.New("unit test with mocks fail"),
	}

	_, err := mockApp.RetrieveSystemServices()
	if err == nil {
		t.Fatal("Error was expected, got nil")
	}
}