package updates

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"os/exec"
)

type SystemUpdateInformation struct {
	HotFixID    string `json:"HotFixID"`
	Description string `json:"Description"`
	InstalledOn string `json:"InstalledOn"`
	InstalledBy string `json:"InstalledBy"`
}

func RetrieveSystemUpdates() ([]SystemUpdateInformation, error) {
	scriptContent := "Get-CimInstance Win32_QuickFixEngineering | Select-Object HotFixID, Description, InstalledOn, InstalledBy | ConvertTo-Json"

	executableCommand := exec.Command("powershell", "-NoProfile", "-Command", scriptContent)

	var outputBuffer bytes.Buffer
	executableCommand.Stdout = &outputBuffer

	executionError := executableCommand.Run()
	if executionError != nil {
		slog.Error("Failed to execute PowerShell command for system updates",
			slog.String("error_details", executionError.Error()),
		)
		return nil, fmt.Errorf("powershell execution error: %w", executionError)
	}

	trimmedBytes := bytes.TrimSpace(outputBuffer.Bytes())
	if len(trimmedBytes) == 0 {
		return []SystemUpdateInformation{}, nil
	}

	if trimmedBytes[0] == '{' {
		var singleUpdate SystemUpdateInformation
		decodingError := json.Unmarshal(trimmedBytes, &singleUpdate)
		if decodingError != nil {
			slog.Error("Failed to decode single JSON update object",
				slog.String("error_details", decodingError.Error()),
			)
			return nil, fmt.Errorf("json parsing error for single update object: %w", decodingError)
		}
		return []SystemUpdateInformation{singleUpdate}, nil
	}

	var updatesList []SystemUpdateInformation
	decodingError := json.Unmarshal(trimmedBytes, &updatesList)
	if decodingError != nil {
		slog.Error("Failed to decode JSON updates list array",
			slog.String("error_details", decodingError.Error()),
		)
		return nil, fmt.Errorf("json parsing error for updates list: %w", decodingError)
	}

	return updatesList, nil
}
