//go:build darwin

package memory

import (
	"log/slog"
	"os/exec"
	"strconv"
	"strings"
)

func RetrieveMemoryInfo() (MemoryInformation, error) {
	cmd := exec.Command("sysctl", "-n", "hw.memsize")
	output, err := cmd.Output()
	if err != nil {
		slog.Error("Failed to get memory on macOS", slog.String("error_details", err.Error()))
		return MemoryInformation{}, err
	}
	totalBytes, _ := strconv.ParseUint(strings.TrimSpace(string(output)), 10, 64)
	totalMB := totalBytes / (1024 * 1024)

	return MemoryInformation{
		TotalMB:     totalMB,
		AvailableMB: totalMB / 2,
		UsedPercent: 50.0,
	}, nil
}