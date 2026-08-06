//go:build unit

package sensors

import (
	"errors"
	"testing"
)

func TestSensors_MockSuccess(t *testing.T) {
	mock := &MockSensors{
		DataToReturn: []SensorInformation{
			{
				Name:  "CPU Zone",
				Value: 50.5,
				Unit:  "°C",
			},
		},
	}

	infos, err := mock.RetrieveSensorsInfo()
	if err != nil {
		t.Fatalf("Absence of an error was expected: %v", err)
	}

	if len(infos) != 1 || infos[0].Name != "CPU Zone" {
		t.Errorf("Unexpected sensors data returned: %v", infos)
	}
}

func TestSensors_MockError(t *testing.T) {
	mock := &MockSensors{
		MockError: errors.New("unit test with mocks fail"),
	}

	_, err := mock.RetrieveSensorsInfo()
	if err == nil {
		t.Fatal("Error was expected, got nil")
	}
}