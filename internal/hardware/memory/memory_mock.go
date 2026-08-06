package memory

type MockMemory struct {
	DataToReturn *MemoryInformation
	MockError    error
}

var _ Memory = (*MockMemory)(nil)

func (mock *MockMemory) RetrieveMemoryInfo() (MemoryInformation, error) {
	if mock.MockError != nil {
		return MemoryInformation{}, mock.MockError
	}

	if mock.DataToReturn == nil {
		return MemoryInformation{
			TotalMB:     12,
			AvailableMB: 13,
			UsedPercent: 1.23,
		}, nil
	}

	return *mock.DataToReturn, nil
}