//go:build darwin

package services

import (
	"bytes"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
)

func (*SoftwareServices)RetrieveSystemServices() ([]ServiceInformation, error) {
	executableCommand := exec.Command("launchctl", "list")

	var outputBuffer bytes.Buffer
	executableCommand.Stdout = &outputBuffer

	executionError := executableCommand.Run()
	if executionError != nil {
		slog.Error("Failed to execute launchctl for services",
			slog.String("error_details", executionError.Error()),
		)
		return nil, fmt.Errorf("launchctl execution error: %w", executionError)
	}

	lines := strings.Split(outputBuffer.String(), "\n")
	var servicesList []ServiceInformation

	// Skip header line
	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 3 {
			servicesList = append(servicesList, ServiceInformation{
				Name:        fields[2],
				DisplayName: fields[2],
				State:       fields[0], // PID or "-" if not running
				StartMode:   "auto",
				StartName:   "system",
			})
		}
	}

	return servicesList, nil
}