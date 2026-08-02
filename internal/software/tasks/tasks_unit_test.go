//go:build unit

package tasks

import (
	"errors"
	"testing"
)

func TestTasks_MockSuccess(t *testing.T) {
	mockApp := &MockTasks{
		DataToReturn: []ScheduledTaskInformation{
			{
				TaskName: "TaskName",
				TaskPath: "TaskPath",
				State:    1,
			},
		},
	}

	_, err := mockApp.RetrieveScheduledTasks()
	if err != nil {
		t.Fatalf("Absence of an error was expected: %v", err)
	}
}

func TestTasks_MockError(t *testing.T) {
	mockApp := &MockTasks{
		MockError: errors.New("unit test with mocks fail"),
	}

	_, err := mockApp.RetrieveScheduledTasks()
	if err == nil {
		t.Fatal("Error was expected, got nil")
	}
}
