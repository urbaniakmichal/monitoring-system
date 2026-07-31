//go:build windows

package tasks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"os/exec"
)

func RetrieveScheduledTasks() ([]ScheduledTaskInformation, error) {
	scriptContent := "Get-ScheduledTask | Select-Object TaskName, TaskPath, State | ConvertTo-Json"

	executableCommand := exec.Command("powershell", "-NoProfile", "-Command", scriptContent)

	var outputBuffer bytes.Buffer
	executableCommand.Stdout = &outputBuffer

	executionError := executableCommand.Run()
	if executionError != nil {
		slog.Error("Failed to execute PowerShell command for scheduled tasks",
			slog.String("error_details", executionError.Error()),
		)
		return nil, fmt.Errorf("powershell execution error: %w", executionError)
	}

	trimmedBytes := bytes.TrimSpace(outputBuffer.Bytes())
	if len(trimmedBytes) == 0 {
		return []ScheduledTaskInformation{}, nil
	}

	if trimmedBytes[0] == '{' {
		var singleScheduledTask ScheduledTaskInformation
		decodingError := json.Unmarshal(trimmedBytes, &singleScheduledTask)
		if decodingError != nil {
			slog.Error("Failed to decode single JSON scheduled task object",
				slog.String("error_details", decodingError.Error()),
			)
			return nil, fmt.Errorf("json parsing error for single scheduled task object: %w", decodingError)
		}
		return []ScheduledTaskInformation{singleScheduledTask}, nil
	}

	var scheduledTasksList []ScheduledTaskInformation
	decodingError := json.Unmarshal(trimmedBytes, &scheduledTasksList)
	if decodingError != nil {
		slog.Error("Failed to decode JSON scheduled tasks list array",
			slog.String("error_details", decodingError.Error()),
		)
		return nil, fmt.Errorf("json parsing error for scheduled tasks list: %w", decodingError)
	}

	return scheduledTasksList, nil
}