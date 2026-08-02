package tasks

type ScheduledTaskInformation struct {
	TaskName string `json:"TaskName"`
	TaskPath string `json:"TaskPath"`
	State    int    `json:"State"`
}

type Tasks interface {
	RetrieveScheduledTasks() ([]ScheduledTaskInformation, error)
}

type SoftwareTasks struct{}

var _ Tasks = (*SoftwareTasks)(nil)