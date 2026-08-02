package startup

type MockStartup struct {
	DataToReturn []StartupCommandInformation
	MockError    error
}

var _ Startup = (*MockStartup)(nil)

func (mock *MockStartup) RetrieveStartupCommands() ([]StartupCommandInformation, error) {
	if mock.MockError != nil {
		return nil, mock.MockError
	}

	if mock.DataToReturn == nil {
		return []StartupCommandInformation{
			{
				Name:     "Name",
				Command:  "Command",
				Location: "Location",
				User:     "User",
			},
		}, nil
	}

	return mock.DataToReturn, nil
}