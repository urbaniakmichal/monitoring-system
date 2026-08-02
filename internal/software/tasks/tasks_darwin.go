//go:build darwin

package tasks

import (
	"bytes"
	"fmt"
	"log/slog"
	"os/exec"
	"strconv"
	"strings"
)

func (*SoftwareTasks)RetrieveScheduledTasks() ([]ScheduledTaskInformation, error) {
	executableCommand := exec.Command("launchctl", "list")

	var outputBuffer bytes.Buffer
	executableCommand.Stdout = &outputBuffer

	executionError := executableCommand.Run()
	if executionError != nil {
		slog.Error("Failed to execute launchctl for scheduled tasks",
			slog.String("error_details", executionError.Error()),
		)
		return nil, fmt.Errorf("launchctl execution error: %w", executionError)
	}

	lines := strings.Split(outputBuffer.String(), "\n")
	var tasksList []ScheduledTaskInformation

	// Skip header line
	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 3 {
			// fields[0] as usually PID
			stateVal := 0
			if fields[0] != "-" {
				if parsedPid, err := strconv.Atoi(fields[0]); err == nil {
					stateVal = parsedPid
				}
			}

			tasksList = append(tasksList, ScheduledTaskInformation{
				TaskName: fields[2],
				TaskPath: "/Library/LaunchAgents",
				State:    stateVal,
			})
		}
	}

	return tasksList, nil
}