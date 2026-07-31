package tasks

type ScheduledTaskInformation struct {
	TaskName string `json:"TaskName"`
	TaskPath string `json:"TaskPath"`
	State    int `json:"State"`
}