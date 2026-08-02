//go:build unit

package network

import (
	"errors"
	"testing"
)

func TestNetworks_MockSuccess(t *testing.T) {
	mockNetwork := &MockNetworks{
		DataToReturn: []NetworkAdapterInformation{
			{
				Caption:     "Caption",
				Description: "Description",
				IPAddress:   "IPAddress",
				MACAddress:  "MACAddress",
				DHCPEnabled: "DHCPEnabled ",
			},
		},
	}

	_, err := mockNetwork.RetrieveActiveNetworkAdapters()
	if err != nil {
		t.Fatalf("Absence of an error was expected: %v", err)
	}
}

func TestNetworks_MockError(t *testing.T) {
	mockNetwork := &MockNetworks{
		MockError: errors.New("unit test with mocks fail"),
	}

	_, err := mockNetwork.RetrieveActiveNetworkAdapters()
	if err == nil {
		t.Fatal("Error was expected, got nil")
	}
}