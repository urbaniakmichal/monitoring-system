package network

type MockNetworks struct {
	DataToReturn []NetworkAdapterInformation
	MockError    error
}

var _ Networks = (*MockNetworks)(nil)

func (mock *MockNetworks) RetrieveActiveNetworkAdapters() ([]NetworkAdapterInformation, error) {
	if mock.MockError != nil {
		return nil, mock.MockError
	}

	if mock.DataToReturn == nil {
		return []NetworkAdapterInformation{
			{
				Caption:     "Caption",
				Description: "Description",
				IPAddress:   "IPAddress",
				MACAddress:  "MACAddress",
				DHCPEnabled: "DHCPEnabled ",
			},
		}, nil
	}

	return mock.DataToReturn, nil
}