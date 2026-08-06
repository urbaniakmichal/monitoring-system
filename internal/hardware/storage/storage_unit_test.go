//go:build unit

package storage

import (
	"errors"
	"testing"
)

func TestStorage_MockSuccess(t *testing.T) {
	mock := &MockStorage{
		DataToReturn: []StorageInformation{
			{
				Device:      "C:",
				Path:        "C:\\",
				TotalMB:     100000,
				FreeMB:      40000,
				UsedPercent: 60.0,
			},
		},
	}

	infos, err := mock.RetrieveStorageInformation()
	if err != nil {
		t.Fatalf("Absence of an error was expected: %v", err)
	}

	if len(infos) != 1 || infos[0].Device != "C:" {
		t.Errorf("Unexpected storage data returned: %v", infos)
	}
}

func TestStorage_MockError(t *testing.T) {
	mock := &MockStorage{
		MockError: errors.New("unit test with mocks fail"),
	}

	_, err := mock.RetrieveStorageInformation()
	if err == nil {
		t.Fatal("Error was expected, got nil")
	}
}