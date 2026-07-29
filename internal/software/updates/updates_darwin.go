//go:build darwin

package updates

import (
	"bytes"
	"log/slog"
	"os/exec"
	"strings"
)

type SystemUpdateInformation struct {
	HotFixID    string `json:"HotFixID"`
	Description string `json:"Description"`
	InstalledOn string `json:"InstalledOn"`
	InstalledBy string `json:"InstalledBy"`
}

func RetrieveSystemUpdates() ([]SystemUpdateInformation, error) {
	executableCommand := exec.Command("softwareupdate", "--history")

	var outputBuffer bytes.Buffer
	executableCommand.Stdout = &outputBuffer

	executionError := executableCommand.Run()
	if executionError != nil {
		slog.Error("Failed to execute softwareupdate for system updates",
			slog.String("error_details", executionError.Error()),
		)
		// Return empty slice gracefully if history is unavailable or empty
		return []SystemUpdateInformation{}, nil
	}

	lines := strings.Split(outputBuffer.String(), "\n")
	var updatesList []SystemUpdateInformation

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "Display Name") || strings.HasPrefix(trimmed, "---") {
			continue
		}
		parts := strings.Split(trimmed, ",")
		if len(parts) >= 2 {
			updatesList = append(updatesList, SystemUpdateInformation{
				HotFixID:    strings.TrimSpace(parts[0]),
				Description: strings.TrimSpace(parts[0]),
				InstalledOn: strings.TrimSpace(parts[len(parts)-1]),
				InstalledBy: "Apple Software Update",
			})
		} else if len(parts) == 1 && trimmed != "" {
			updatesList = append(updatesList, SystemUpdateInformation{
				HotFixID:    trimmed,
				Description: trimmed,
				InstalledOn: "",
				InstalledBy: "Apple Software Update",
			})
		}
	}

	return updatesList, nil
}