package applications

type MockApplications struct {
	DataToReturn []ApplicationInformation
	MockError    error
}

var _ Applications = (*MockApplications)(nil)

func (mock *MockApplications) RetrieveInstalledApplications() ([]ApplicationInformation, error) {
	if mock.MockError != nil {
		return nil, mock.MockError
	}

	if mock.DataToReturn == nil {
		return []ApplicationInformation{
			{
				DisplayName:    "Mock display name",
				DisplayVersion: "Mock display version",
				Publisher:      "Mock Publisher",
				InstallDate:    "20/02/2000",
			},
		}, nil
	}

	return mock.DataToReturn, nil
}