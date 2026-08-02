package tasks

type MockTasks struct {
	DataToReturn []ScheduledTaskInformation
	MockError    error
}

var _ Tasks = (*MockTasks)(nil)

func (mock *MockTasks) RetrieveScheduledTasks() ([]ScheduledTaskInformation, error) {
	if mock.MockError != nil {
		return nil, mock.MockError
	}

	if mock.DataToReturn == nil {
		return []ScheduledTaskInformation{
			{
				TaskName: "TaskName",
				TaskPath: "TaskPath",
				State:    1,
			},
		}, nil
	}

	return mock.DataToReturn, nil
}