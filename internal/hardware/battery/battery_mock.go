package battery

type MockBattery struct {
	DataToReturn []BatteryInformation
	MockError    error
}

var _ Battery = (*MockBattery)(nil)

func (mock *MockBattery) RetrieveBatteryInfo() ([]BatteryInformation, error) {
	if mock.MockError != nil {
		return nil, mock.MockError
	}

	if mock.DataToReturn == nil {
		return []BatteryInformation{
			{
				Name:       "Name",
				Percent:    12.34,
				IsCharging: true,
				Status:     "Status",
			},
		}, nil
	}

	return mock.DataToReturn, nil
}