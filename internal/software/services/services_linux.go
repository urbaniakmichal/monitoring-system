//go:build linux

package services

import (
	"bytes"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
)

func RetrieveSystemServices() ([]ServiceInformation, error) {
	cmd := exec.Command("systemctl", "list-units", "--type=service", "--no-legend", "--all")
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		slog.Error("Failed to execute systemctl", slog.String("error_details", err.Error()))
		return nil, fmt.Errorf("systemctl error: %w", err)
	}

	var services []ServiceInformation
	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 4 {
			name := fields[0]
			state := fields[2]
			services = append(services, ServiceInformation{
				Name:        name,
				DisplayName: name,
				State:       state,
				StartMode:   "enabled",
				StartName:   "systemd",
			})
		}
	}
	return services, nil
}