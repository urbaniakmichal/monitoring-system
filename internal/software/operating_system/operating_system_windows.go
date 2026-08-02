//go:build windows

package operating_system

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"os/exec"
)

func (*SoftwareOperatingSystem)RetrieveOperatingSystemInformation() (OperatingSystemInformation, error) {
	scriptContent := "Get-CimInstance Win32_OperatingSystem | Select-Object Caption, Version, BuildNumber, Manufacturer, OSArchitecture, InstallDate, LastBootUpTime | ConvertTo-Json"

	executableCommand := exec.Command("powershell", "-NoProfile", "-Command", scriptContent)

	var outputBuffer bytes.Buffer
	executableCommand.Stdout = &outputBuffer

	executionError := executableCommand.Run()
	if executionError != nil {
		slog.Error("Failed to execute PowerShell command for operating system information",
			slog.String("error_details", executionError.Error()),
		)
		return OperatingSystemInformation{}, fmt.Errorf("powershell execution error: %w", executionError)
	}

	trimmedBytes := bytes.TrimSpace(outputBuffer.Bytes())
	if len(trimmedBytes) == 0 {
		return OperatingSystemInformation{}, fmt.Errorf("powershell output is empty")
	}

	var systemInformation OperatingSystemInformation
	decodingError := json.Unmarshal(trimmedBytes, &systemInformation)
	if decodingError != nil {
		slog.Error("Failed to decode JSON operating system object",
			slog.String("error_details", decodingError.Error()),
		)
		return OperatingSystemInformation{}, fmt.Errorf("json parsing error for operating system information: %w", decodingError)
	}

	return systemInformation, nil
}