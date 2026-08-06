package storage

type MockStorage struct {
	DataToReturn []StorageInformation
	MockError    error
}

var _ Storage = (*MockStorage)(nil)

func (mock *MockStorage) RetrieveStorageInformation() ([]StorageInformation, error) {
	if mock.MockError != nil {
		return nil, mock.MockError
	}

	if mock.DataToReturn == nil {
		return []StorageInformation{
			{
				Device:      "/dev/sda1",
				Path:        "/",
				TotalMB:     500000,
				FreeMB:      250000,
				UsedPercent: 50.0,
			},
		}, nil
	}

	return mock.DataToReturn, nil
}