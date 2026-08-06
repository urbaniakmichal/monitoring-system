package gpu

type MockGpu struct {
	DataToReturn []GPUInformation
	MockError    error
}

var _ Gpu = (*MockGpu)(nil)

func (mock *MockGpu) RetrieveGPUInfo() ([]GPUInformation, error) {
	if mock.MockError != nil {
		return nil, mock.MockError
	}

	if mock.DataToReturn == nil {
		return []GPUInformation{
			{
				Name:   "Name",
				Vendor: "Vendor",
			},
		}, nil
	}

	return mock.DataToReturn, nil
}