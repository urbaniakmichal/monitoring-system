//go:build windows

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

	cmd := exec.Command("powershell", "-Command", "(Get-CimInstance Win32_Processor).Name")
	output, err := cmd.Output()
	if err != nil {
		slog.Error("Failed to retrieve CPU name on Windows", slog.String("error_details", err.Error()))
		info.ModelName = "Unknown Windows CPU"
		return info, nil
	}

	info.ModelName = strings.TrimSpace(string(output))
	return info, nil
}