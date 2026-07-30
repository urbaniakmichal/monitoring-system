//go:build darwin

package cpu

import (
	"log/slog"
	"os/exec"
	"runtime"
	"strings"
)

func RetrieveCPUInfo() (CPUInformation, error) {
	info := CPUInformation{
		Cores: runtime.NumCPU(),
	}

	cmd := exec.Command("sysctl", "-n", "machdep.cpu.brand_string")
	output, err := cmd.Output()
	if err != nil {
		slog.Error("Failed to retrieve CPU name on macOS", slog.String("error_details", err.Error()))
		info.ModelName = "Unknown Mac CPU"
		return info, nil
	}

	info.ModelName = strings.TrimSpace(string(output))
	return info, nil
}