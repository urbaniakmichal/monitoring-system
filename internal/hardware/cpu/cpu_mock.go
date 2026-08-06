package cpu

type MockCpuInformation struct {
	DataToReturn *CPUInformation
	MockError    error
}

var _ Cpu = (*MockCpuInformation)(nil)

func (mock *MockCpuInformation) RetrieveCPUInfo() (CPUInformation, error) {
	if mock.MockError != nil {
		return CPUInformation{}, mock.MockError
	}

	if mock.DataToReturn == nil {
		return CPUInformation{
			ModelName:    "ModelName",
			Cores:        123,
			UsagePercent: 1.23,
		}, nil
	}

	return *mock.DataToReturn, nil
}