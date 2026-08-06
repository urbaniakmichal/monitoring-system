//go:build windows

package gpu

import (
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
)

func (*HardwareGpu)RetrieveGPUInfo() ([]GPUInformation, error) {
	cmd := exec.Command("powershell", "-Command", "(Get-CimInstance Win32_VideoController).Name")
	output, err := cmd.Output()
	if err != nil {
		slog.Error("Failed to retrieve GPU on Windows", slog.String("error_details", err.Error()))
		return nil, fmt.Errorf("failed to get gpu info: %w", err)
	}

	var gpus []GPUInformation
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			gpus = append(gpus, GPUInformation{
				Name:   trimmed,
				Vendor: "Microsoft WMI",
			})
		}
	}
	return gpus, nil
}