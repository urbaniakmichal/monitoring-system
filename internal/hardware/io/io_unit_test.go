//go:build unit

package io

import (
	"errors"
	"testing"
)

func TestIo_MockSuccess(t *testing.T) {
	mock := &MockIOStats{
   		DataToReturn: &IOStatistics{
			Network: map[string]NetworkIO{
				"eth0": {
					BytesSent: 1024,
					BytesRecv: 2048,
				},
			},
			Disk: map[string]DiskIO{
				"sda": {
					ReadBytes:  5120,
					WriteBytes: 10240,
				},
			},
		},
	}

	_, err := mock.RetrieveIOStatistics()
	if err != nil {
		t.Fatalf("Absence of an error was expected: %v", err)
	}
}

func TestIo_MockError(t *testing.T) {
	mock := &MockIOStats{
		MockError: errors.New("unit test with mocks fail"),
	}

	_, err := mock.RetrieveIOStatistics()
	if err == nil {
		t.Fatal("Error was expected, got nil")
	}
}