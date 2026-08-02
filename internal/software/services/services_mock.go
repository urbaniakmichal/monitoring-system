package services

type MockServices struct {
	DataToReturn []ServiceInformation
	MockError    error
}

var _ Services = (*MockServices)(nil)

func (mock *MockServices) RetrieveSystemServices() ([]ServiceInformation, error) {
	if mock.MockError != nil {
		return nil, mock.MockError
	}

	if mock.DataToReturn == nil {
		return []ServiceInformation{
			{
				Name:        "Name",
				DisplayName: "DisplayName",
				State:       "State",
				StartMode:   "StartMode",
				StartName:   "StartName",
			},
		}, nil
	}

	return mock.DataToReturn, nil
}