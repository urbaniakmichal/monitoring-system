//go:build unit

package users

import (
	"errors"
	"testing"
)


func TestUsers_MockSuccess(t *testing.T) {
	mockApp := &MockUsers{
		DataToReturn: []UserInformation{
			{
				Username:    "HotFixID",
				Terminal: "Description",
				Host: "InstalledOn",
			},
		},
	}

	_, err := mockApp.RetrieveUsersInfo()
	if err != nil {
		t.Fatalf("Absence of an error was expected: %v", err)
	}
}

func TestUsers_MockError(t *testing.T) {
	mockApp := &MockUsers{
		MockError: errors.New("unit test with mocks fail"),
	}

	_, err := mockApp.RetrieveUsersInfo()
	if err == nil {
		t.Fatal("Error was expected, got nil")
	}
}
