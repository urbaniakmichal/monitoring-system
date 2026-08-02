package operating_system

type MockOperatingSystem struct {
	DataToReturn *OperatingSystemInformation
	MockError    error
}

var _ OperatingSystem = (*MockOperatingSystem)(nil)

func (mock *MockOperatingSystem) RetrieveOperatingSystemInformation() (OperatingSystemInformation, error) {
	if mock.MockError != nil {
		return OperatingSystemInformation{}, mock.MockError
	}

	if mock.DataToReturn == nil {
		return OperatingSystemInformation{
			Caption:        "Caption",
			Version:        "Version",
			BuildNumber:    "BuildNumber",
			Manufacturer:   "Manufacturer",
			OSArchitecture: "OSArchitecture",
			InstallDate:    "InstallDate",
			LastBootTime:   "LastBootUpTime",
		}, nil
	}

	return *mock.DataToReturn, nil
}