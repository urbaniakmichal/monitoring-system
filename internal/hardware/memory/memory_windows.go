//go:build windows

package memory

import (
	"log/slog"
	"os/exec"
	"strconv"
	"strings"
)

func (*HardwareMemory)RetrieveMemoryInfo() (MemoryInformation, error) {
	cmdTotal := exec.Command("powershell", "-Command", "(Get-CimInstance Win32_OperatingSystem).TotalVisibleMemorySize")
	outTotal, err := cmdTotal.Output()
	if err != nil {
		slog.Error("Failed to get total memory on Windows", slog.String("error_details", err.Error()))
		return MemoryInformation{}, err
	}
	totalKB, _ := strconv.ParseUint(strings.TrimSpace(string(outTotal)), 10, 64)

	cmdFree := exec.Command("powershell", "-Command", "(Get-CimInstance Win32_OperatingSystem).FreePhysicalMemory")
	outFree, err := cmdFree.Output()
	if err != nil {
		slog.Error("Failed to get free memory on Windows", slog.String("error_details", err.Error()))
		return MemoryInformation{}, err
	}
	freeKB, _ := strconv.ParseUint(strings.TrimSpace(string(outFree)), 10, 64)

	totalMB := totalKB / 1024
	freeMB := freeKB / 1024
	usedMB := totalMB - freeMB
	var usedPercent float64
	if totalMB > 0 {
		usedPercent = (float64(usedMB) / float64(totalMB)) * 100
	}

	return MemoryInformation{
		TotalMB:     totalMB,
		AvailableMB: freeMB,
		UsedPercent: usedPercent,
	}, nil
}