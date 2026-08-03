package users

type MockUsers struct {
	DataToReturn []UserInformation
	MockError    error
}

var _ User = (*MockUsers)(nil)

func (mock *MockUsers) RetrieveUsersInfo() ([]UserInformation, error) {
	if mock.MockError != nil {
		return nil, mock.MockError
	}

	if mock.DataToReturn == nil {
		return []UserInformation{
			{
				Username: "Username",
				Terminal: "Terminal",
				Host:     "Host",
			},
		}, nil
	}

	return mock.DataToReturn, nil
}