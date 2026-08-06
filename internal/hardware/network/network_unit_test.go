//go:build unit

package network

import (
	"errors"
	"testing"
)

func TestMockNetwork_DefaultData(t *testing.T) {
	mock := &MockNetwork{}

	infos, err := mock.RetrieveNetworkInfo()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if len(infos) != 1 {
		t.Fatalf("expected 1 interface, got: %d", len(infos))
	}

	if infos[0].Name != "eth0" {
		t.Errorf("expected name 'eth0', got: '%s'", infos[0].Name)
	}
}

func TestMockNetwork_CustomData(t *testing.T) {
	customData := []NetworkInformation{
		{
			Name:         "wlan0",
			HardwareAddr: "aa:bb:cc:dd:ee:ff",
			IPAddresses:  []string{"10.0.0.5"},
			Flags:        "up",
		},
	}

	mock := &MockNetwork{
		DataToReturn: customData,
	}

	infos, err := mock.RetrieveNetworkInfo()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if len(infos) != 1 || infos[0].Name != "wlan0" {
		t.Errorf("did not return expected custom data: %v", infos)
	}
}

func TestMockNetwork_Error(t *testing.T) {
	mock := &MockNetwork{
		MockError: errors.New("network error"),
	}

	_, err := mock.RetrieveNetworkInfo()
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != "network error" {
		t.Errorf("expected message 'network error', got: '%v'", err)
	}
}