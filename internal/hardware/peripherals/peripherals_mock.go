package peripherals

type MockPeripherals struct {
	DataToReturn *PeripheralsInformation
	MockError    error
}

var _ Peripherals = (*MockPeripherals)(nil)

func (mock *MockPeripherals) RetrievePeripheralsInfo() (PeripheralsInformation, error) {
	if mock.MockError != nil {
		return PeripheralsInformation{}, mock.MockError
	}

	if mock.DataToReturn == nil {
		return PeripheralsInformation{
			Displays: []DisplayInformation{
				{
					Name:       "Display1",
					Resolution: "1920x1080",
					IsPrimary:  true,
				},
			},
			USB: []USBDeviceInformation{
				{
					Name:         "USB Keyboard",
					Manufacturer: "Generic",
					DeviceID:     "USB\\VID_0000&PID_0000",
				},
			},
		}, nil
	}

	return *mock.DataToReturn, nil
}