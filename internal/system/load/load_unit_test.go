//go:build unit

package load

import (
	"errors"
	"testing"
)


func TestSystemLoad_MockSuccess(t *testing.T) {
	mock := &MockSystemLoadStruct{
   		DataToReturn: &LoadInformation{
			Load1:     1.1,
			Load5:           1.2,
			Load15: 1.3,
		},
	}

	_, err := mock.RetrieveLoadInfo()
	if err != nil {
		t.Fatalf("Absence of an error was expected: %v", err)
	}
}

func TestSystemLoad_MockError(t *testing.T) {
	mock := &MockSystemLoadStruct{
		MockError: errors.New("unit test with mocks fail"),
	}

	_, err := mock.RetrieveLoadInfo()
	if err == nil {
		t.Fatal("Error was expected, got nil")
	}
}