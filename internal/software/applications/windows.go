package applications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"os/exec"
)

type ApplicationInformation struct {
	DisplayName    string `json:"DisplayName"`
	DisplayVersion string `json:"DisplayVersion"`
	Publisher      string `json:"Publisher"`
	InstallDate    string `json:"InstallDate"`
}

func RetrieveInstalledApplications() ([]ApplicationInformation, error) {
	scriptContent := "$paths = @('HKLM:\\Software\\Microsoft\\Windows\\CurrentVersion\\Uninstall\\*', 'HKLM:\\Software\\Wow6432Node\\Microsoft\\Windows\\CurrentVersion\\Uninstall\\*'); Get-ItemProperty -Path $paths -ErrorAction SilentlyContinue | Where-Object { $_.DisplayName -ne $null } | Select-Object DisplayName, DisplayVersion, Publisher, @{Name='InstallDate';Expression={ \"$($_.InstallDate)\" }} | Sort-Object DisplayName -Unique | ConvertTo-Json"

	executableCommand := exec.Command("powershell", "-NoProfile", "-Command", scriptContent)

	var outputBuffer bytes.Buffer
	executableCommand.Stdout = &outputBuffer

	executionError := executableCommand.Run()
	if executionError != nil {
		slog.Error("Failed to execute PowerShell command for installed applications",
			slog.String("error_details", executionError.Error()),
		)
		return nil, fmt.Errorf("powershell execution error: %w", executionError)
	}

	trimmedBytes := bytes.TrimSpace(outputBuffer.Bytes())
	if len(trimmedBytes) == 0 {
		return []ApplicationInformation{}, nil
	}

	if trimmedBytes[0] == '{' {
		var singleApplication ApplicationInformation
		decodingError := json.Unmarshal(trimmedBytes, &singleApplication)
		if decodingError != nil {
			slog.Error("Failed to decode single JSON application object",
				slog.String("error_details", decodingError.Error()),
			)
			return nil, fmt.Errorf("json parsing error for single application object: %w", decodingError)
		}
		return []ApplicationInformation{singleApplication}, nil
	}

	var applicationsList []ApplicationInformation
	decodingError := json.Unmarshal(trimmedBytes, &applicationsList)
	if decodingError != nil {
		slog.Error("Failed to decode JSON applications list array",
			slog.String("error_details", decodingError.Error()),
		)
		return nil, fmt.Errorf("json parsing error for applications list: %w", decodingError)
	}

	return applicationsList, nil
}
