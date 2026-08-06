package io

type MockIOStats struct {
	DataToReturn *IOStatistics
	MockError    error
}

func (mock *MockIOStats) RetrieveIOStatistics() (IOStatistics, error) {
	if mock.MockError != nil {
		return IOStatistics{}, mock.MockError
	}

	if mock.DataToReturn == nil {
		return IOStatistics{
			Network: map[string]NetworkIO{
				"eth0": {
					BytesSent: 1024,
					BytesRecv: 2048,
				},
			},
			Disk: map[string]DiskIO{
				"sda": {
					ReadBytes:  5120,
					WriteBytes: 10240,
				},
			},
		}, nil
	}

	return *mock.DataToReturn, nil
}