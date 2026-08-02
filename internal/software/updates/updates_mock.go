package updates

type MockUpdates struct {
	DataToReturn []SystemUpdateInformation
	MockError    error
}

var _ Updates = (*MockUpdates)(nil)

func (mock *MockUpdates) RetrieveSystemUpdates() ([]SystemUpdateInformation, error) {
	if mock.MockError != nil {
		return nil, mock.MockError
	}

	if mock.DataToReturn == nil {
		return []SystemUpdateInformation{
			{
				HotFixID:    "HotFixID",
				Description: "Description",
				InstalledOn: "InstalledOn",
				InstalledBy: "InstalledBy",
			},
		}, nil
	}

	return mock.DataToReturn, nil
}