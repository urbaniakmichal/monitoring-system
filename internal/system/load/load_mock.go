package load

type MockSystemLoadStruct struct {
	DataToReturn *LoadInformation
	MockError    error
}

var _ SystemLoad = (*MockSystemLoadStruct)(nil)

func (mock *MockSystemLoadStruct) RetrieveLoadInfo() (LoadInformation, error) {
	if mock.MockError != nil {
		return LoadInformation{}, mock.MockError
	}

	if mock.DataToReturn == nil {
		return LoadInformation{
			Load1:  1.1,
			Load5:  1.2,
			Load15: 1.3,
		}, nil
	}

	return *mock.DataToReturn, nil
}