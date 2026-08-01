package drivers

type MockDrivers struct {
	AppsToReturn []DriverInformation
	MockError    error
}

var _ Drivers = (*MockDrivers)(nil)

func (mock *MockDrivers) RetrieveInstalledDrivers() ([]DriverInformation, error) {
	if mock.MockError != nil {
		return nil, mock.MockError
	}

	if mock.AppsToReturn == nil {
		return []DriverInformation{
			{
				DeviceName:    "DeviceName",
				DriverVersion: "DriverVersion",
				Manufacturer:  "Manufacturer",
				DriverName:    "DriverName",
			},
		}, nil
	}

	return mock.AppsToReturn, nil
}