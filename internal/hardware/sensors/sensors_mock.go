package sensors

type MockSensors struct {
	DataToReturn []SensorInformation
	MockError    error
}

var _ Sensors = (*MockSensors)(nil)

func (mock *MockSensors) RetrieveSensorsInfo() ([]SensorInformation, error) {
	if mock.MockError != nil {
		return nil, mock.MockError
	}

	if mock.DataToReturn == nil {
		return []SensorInformation{
			{
				Name:  "Thermal Zone 0",
				Value: 45.0,
				Unit:  "°C",
			},
		}, nil
	}

	return mock.DataToReturn, nil
}