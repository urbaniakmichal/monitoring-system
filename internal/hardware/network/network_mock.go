package network

type MockNetwork struct {
	DataToReturn []NetworkInformation
	MockError    error
}

var _ Network = (*MockNetwork)(nil)

func (mock *MockNetwork) RetrieveNetworkInfo() ([]NetworkInformation, error) {
	if mock.MockError != nil {
		return nil, mock.MockError
	}

	if mock.DataToReturn == nil {
		return []NetworkInformation{
			{
				Name:         "eth0",
				HardwareAddr: "00:11:22:33:44:55",
				IPAddresses:  []string{"192.168.1.10"},
				Flags:        "up,broadcast,running",
			},
		}, nil
	}

	return mock.DataToReturn, nil
}