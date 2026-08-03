package system

type MockSystemInformation struct {
	DataToReturn *SystemInformation
	MockError    error
}

var _ System = (*MockSystemInformation)(nil)

func (mock *MockSystemInformation) RetrieveSystemInfo() (SystemInformation, error) {
	if mock.MockError != nil {
		return SystemInformation{}, mock.MockError
	}

	if mock.DataToReturn == nil {
		return SystemInformation{
			Hostname:     "Hostname",
			OS:           "OS",
			Architecture: "Architecture",
			Uptime:       "Uptime",
			TopProcesses: []ProcessInformation{
				{
					PID:        123,
					Name:       "Name",
					CPUPercent: 1.12,
					MemoryMB:   666,
				},
			},
		}, nil
	}

	return *mock.DataToReturn, nil
}