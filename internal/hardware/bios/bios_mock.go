package bios

type MockBios struct {
	DataToReturn *BiosInformation
	MockError    error
}

var _ Bios = (*MockBios)(nil)

func (mock *MockBios) RetrieveBiosInformation() (BiosInformation, error) {
	if mock.MockError != nil {
		return BiosInformation{}, mock.MockError
	}

	if mock.DataToReturn == nil {
		return BiosInformation{
			Manufacturer: "Manufacturer",
			Version:      "Version",
			ReleaseDate:  "ReleaseDate",
			SerialNumber: "SerialNumber",
		}, nil
	}

	return *mock.DataToReturn, nil
}