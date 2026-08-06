//go:build unit

package peripherals

import (
	"errors"
	"testing"
)

func TestPeripherals_MockSuccess(t *testing.T) {
	mock := &MockPeripherals{
		DataToReturn: &PeripheralsInformation{
			Displays: []DisplayInformation{
				{
					Name:       "TestDisplay",
					Resolution: "2560x1440",
					IsPrimary:  true,
				},
			},
		},
	}

	info, err := mock.RetrievePeripheralsInfo()
	if err != nil {
		t.Fatalf("Absence of an error was expected: %v", err)
	}

	if len(info.Displays) != 1 || info.Displays[0].Name != "TestDisplay" {
		t.Errorf("Unexpected peripherals data returned: %v", info)
	}
}

func TestPeripherals_MockError(t *testing.T) {
	mock := &MockPeripherals{
		MockError: errors.New("unit test with mocks fail"),
	}

	_, err := mock.RetrievePeripheralsInfo()
	if err == nil {
		t.Fatal("Error was expected, got nil")
	}
}