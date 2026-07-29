//go:build windows

package drivers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"os/exec"
)

type DriverInformation struct {
	Name        string `json:"Name"`
	DisplayName string `json:"DisplayName"`
	State       string `json:"State"`
	Status      string `json:"Status"`
	StartMode   string `json:"StartMode"`
}

func RetrieveInstalledDrivers() ([]DriverInformation, error) {
	scriptContent := "Get-CimInstance Win32_SystemDriver | Select-Object Name, DisplayName, State, Status, StartMode | ConvertTo-Json"

	executableCommand := exec.Command("powershell", "-NoProfile", "-Command", scriptContent)

	var outputBuffer bytes.Buffer
	executableCommand.Stdout = &outputBuffer

	executionError := executableCommand.Run()
	if executionError != nil {
		slog.Error("Failed to execute PowerShell command for system drivers",
			slog.String("error_details", executionError.Error()),
		)
		return nil, fmt.Errorf("powershell execution error: %w", executionError)
	}

	trimmedBytes := bytes.TrimSpace(outputBuffer.Bytes())
	if len(trimmedBytes) == 0 {
		return []DriverInformation{}, nil
	}

	if trimmedBytes[0] == '{' {
		var singleDriver DriverInformation
		decodingError := json.Unmarshal(trimmedBytes, &singleDriver)
		if decodingError != nil {
			slog.Error("Failed to decode single JSON driver object",
				slog.String("error_details", decodingError.Error()),
			)
			return nil, fmt.Errorf("json parsing error for single driver object: %w", decodingError)
		}
		return []DriverInformation{singleDriver}, nil
	}

	var driversList []DriverInformation
	decodingError := json.Unmarshal(trimmedBytes, &driversList)
	if decodingError != nil {
		slog.Error("Failed to decode JSON drivers list array",
			slog.String("error_details", decodingError.Error()),
		)
		return nil, fmt.Errorf("json parsing error for drivers list: %w", decodingError)
	}

	return driversList, nil
}
