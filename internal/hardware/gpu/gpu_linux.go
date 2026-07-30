//go:build linux

package gpu

import (
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
)

func RetrieveGPUInfo() ([]GPUInformation, error) {
	cmd := exec.Command("lspci")
	output, err := cmd.Output()
	if err != nil {
		slog.Error("Failed to run lspci", slog.String("error_details", err.Error()))
		return nil, fmt.Errorf("failed to get gpu info: %w", err)
	}

	var gpus []GPUInformation
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		lower := strings.ToLower(line)
		if strings.Contains(lower, "vga compatible controller") || strings.Contains(lower, "3d controller") || strings.Contains(lower, "display controller") {
			parts := strings.SplitN(line, ": ", 2)
			name := line
			if len(parts) > 1 {
				name = parts[1]
			}
			gpus = append(gpus, GPUInformation{
				Name:   strings.TrimSpace(name),
				Vendor: "PCI Device",
			})
		}
	}
	return gpus, nil
}