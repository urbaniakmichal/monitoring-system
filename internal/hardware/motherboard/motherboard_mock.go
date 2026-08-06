package motherboard

type MockMotherboard struct {
	DataToReturn *MotherboardInformation
	MockError    error
}

var _ Motherboard = (*MockMotherboard)(nil)

func (mock *MockMotherboard) RetrieveMotherboardInfo() (MotherboardInformation, error) {
	if mock.MockError != nil {
		return MotherboardInformation{}, mock.MockError
	}

	if mock.DataToReturn == nil {
		return MotherboardInformation{
			Manufacturer: "Manufacturer",
			Product:      "Product",
			Version:      "Version",
		}, nil
	}

	return *mock.DataToReturn, nil
}