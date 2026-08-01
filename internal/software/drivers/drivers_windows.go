//go:build windows

package drivers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"os/exec"
)

type winDriver struct {
	Name        string `json:"Name"`
	DisplayName string `json:"DisplayName"`
	State       string `json:"State"`
	
	Status      string `json:"Status"`
	StartMode   string `json:"StartMode"`
}

func (*SoftwareDrivers) RetrieveInstalledDrivers() ([]DriverInformation, error) {
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

	var winDrivers []winDriver
	if trimmedBytes[0] == '{' {
		var single winDriver
		if err := json.Unmarshal(trimmedBytes, &single); err != nil {
			slog.Error("Failed to decode single JSON driver object", slog.String("error_details", err.Error()))
			return nil, fmt.Errorf("json parsing error for single driver object: %w", err)
		}
		winDrivers = []winDriver{single}
	} else {
		if err := json.Unmarshal(trimmedBytes, &winDrivers); err != nil {
			slog.Error("Failed to decode JSON drivers list array", slog.String("error_details", err.Error()))
			return nil, fmt.Errorf("json parsing error for drivers list: %w", err)
		}
	}

	var driversList []DriverInformation
	for _, wd := range winDrivers {
		name := wd.DisplayName
		if name == "" {
			name = wd.Name
		}
		driversList = append(driversList, DriverInformation{
			DeviceName:    name,
			DriverVersion: wd.State,
			Manufacturer:  wd.Status,
			DriverName:    wd.Name,
		})
	}

	return driversList, nil
}