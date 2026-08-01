//go:build linux

package tasks

import (
	"bytes"
	"log/slog"
	"os/exec"
	"strings"
)

func RetrieveScheduledTasks() ([]ScheduledTaskInformation, error) {
	cmd := exec.Command("systemctl", "list-timers", "--no-legend")
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		slog.Error("Failed to list systemd timers", slog.String("error_details", err.Error()))
		return []ScheduledTaskInformation{}, nil
	}

	var tasks []ScheduledTaskInformation
	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) > 0 {
			timerName := fields[len(fields)-1]
			tasks = append(tasks, ScheduledTaskInformation{
				TaskName: timerName,
				TaskPath: "/etc/systemd/system",
				State:    1, // 1 as active timer
			})
		}
	}
	return tasks, nil
}