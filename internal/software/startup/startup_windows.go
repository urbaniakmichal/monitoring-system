//go:build windows

package startup

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"os/exec"
)

func (*SoftwareStartup)RetrieveStartupCommands() ([]StartupCommandInformation, error) {
	scriptContent := "Get-CimInstance Win32_StartupCommand | Select-Object Name, Command, Location, User | ConvertTo-Json"

	executableCommand := exec.Command("powershell", "-NoProfile", "-Command", scriptContent)

	var outputBuffer bytes.Buffer
	executableCommand.Stdout = &outputBuffer

	executionError := executableCommand.Run()
	if executionError != nil {
		slog.Error("Failed to execute PowerShell command for startup commands",
			slog.String("error_details", executionError.Error()),
		)
		return nil, fmt.Errorf("powershell execution error: %w", executionError)
	}

	trimmedBytes := bytes.TrimSpace(outputBuffer.Bytes())
	if len(trimmedBytes) == 0 {
		return []StartupCommandInformation{}, nil
	}

	if trimmedBytes[0] == '{' {
		var singleStartupCommand StartupCommandInformation
		decodingError := json.Unmarshal(trimmedBytes, &singleStartupCommand)
		if decodingError != nil {
			slog.Error("Failed to decode single JSON startup command object",
				slog.String("error_details", decodingError.Error()),
			)
			return nil, fmt.Errorf("json parsing error for single startup command object: %w", decodingError)
		}
		return []StartupCommandInformation{singleStartupCommand}, nil
	}
	var startupCommandsList []StartupCommandInformation
	decodingError := json.Unmarshal(trimmedBytes, &startupCommandsList)
	if decodingError != nil {
		slog.Error("Failed to decode JSON startup commands list array",
			slog.String("error_details", decodingError.Error()),
		)
		return nil, fmt.Errorf("json parsing error for startup commands list: %w", decodingError)
	}

	return startupCommandsList, nil
}