//go:build darwin

package gpu

import (
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
)

func RetrieveGPUInfo() ([]GPUInformation, error) {
	cmd := exec.Command("system_profiler", "SPDisplaysDataType")
	output, err := cmd.Output()
	if err != nil {
		slog.Error("Failed to retrieve GPU on macOS", slog.String("error_details", err.Error()))
		return nil, fmt.Errorf("failed to get gpu info: %w", err)
	}

	var gpus []GPUInformation
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Chipset Model:") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				gpus = append(gpus, GPUInformation{
					Name:   strings.TrimSpace(parts[1]),
					Vendor: "Apple",
				})
			}
		}
	}
	return gpus, nil
}