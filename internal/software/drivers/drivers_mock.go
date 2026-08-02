package drivers

type MockDrivers struct {
	DataToReturn []DriverInformation
	MockError    error
}

var _ Drivers = (*MockDrivers)(nil)

func (mock *MockDrivers) RetrieveInstalledDrivers() ([]DriverInformation, error) {
	if mock.MockError != nil {
		return nil, mock.MockError
	}

	if mock.DataToReturn == nil {
		return []DriverInformation{
			{
				DeviceName:    "DeviceName",
				DriverVersion: "DriverVersion",
				Manufacturer:  "Manufacturer",
				DriverName:    "DriverName",
			},
		}, nil
	}

	return mock.DataToReturn, nil
}