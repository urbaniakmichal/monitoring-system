//go:build windows

package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"os/exec"
)

func RetrieveSystemServices() ([]ServiceInformation, error) {
	scriptContent := "Get-CimInstance Win32_Service | Select-Object Name, DisplayName, State, StartMode, StartName | ConvertTo-Json"

	executableCommand := exec.Command("powershell", "-NoProfile", "-Command", scriptContent)

	var outputBuffer bytes.Buffer
	executableCommand.Stdout = &outputBuffer

	executionError := executableCommand.Run()
	if executionError != nil {
		slog.Error("Failed to execute PowerShell command for system services",
			slog.String("error_details", executionError.Error()),
		)
		return nil, fmt.Errorf("powershell execution error: %w", executionError)
	}

	trimmedBytes := bytes.TrimSpace(outputBuffer.Bytes())
	if len(trimmedBytes) == 0 {
		return []ServiceInformation{}, nil
	}

	if trimmedBytes[0] == '{' {
		var singleService ServiceInformation
		decodingError := json.Unmarshal(trimmedBytes, &singleService)
		if decodingError != nil {
			slog.Error("Failed to decode single JSON service object",
				slog.String("error_details", decodingError.Error()),
			)
			return nil, fmt.Errorf("json parsing error for single service object: %w", decodingError)
		}
		return []ServiceInformation{singleService}, nil
	}

	var servicesList []ServiceInformation
	decodingError := json.Unmarshal(trimmedBytes, &servicesList)
	if decodingError != nil {
		slog.Error("Failed to decode JSON services list array",
			slog.String("error_details", decodingError.Error()),
		)
		return nil, fmt.Errorf("json parsing error for services list: %w", decodingError)
	}

	return servicesList, nil
}