//go:build unit

package battery

import (
	"errors"
	"testing"
)

func TestBattery_MockSuccess(t *testing.T) {
	mockApp := &MockBattery{
		DataToReturn: []BatteryInformation{
			{
				Name:       "Name",
				Percent:    12.34,
				IsCharging: true,
				Status:     "Status",
			},
		},
	}

	_, err := mockApp.RetrieveBatteryInfo()
	if err != nil {
		t.Fatalf("Absence of an error was expected: %v", err)
	}
}

func TestBattery_MockError(t *testing.T) {
	mockApp := &MockBattery{
		MockError: errors.New("unit test with mocks fail"),
	}

	_, err := mockApp.RetrieveBatteryInfo()
	if err == nil {
		t.Fatal("Error was expected, got nil")
	}
}