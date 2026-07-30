//go:build windows

package storage

import (
	"fmt"
	"log/slog"
	"os/exec"
	"strconv"
	"strings"
)

func RetrieveStorageInfo(driveLetter string) (StorageInformation, error) {
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf("$d = Get-CimInstance Win32_LogicalDisk -Filter \"DeviceID='%s:'\"; \"$($d.Size),$($d.FreeSpace)\"", driveLetter))
	output, err := cmd.Output()
	if err != nil {
		slog.Error("Failed to get storage on Windows", slog.String("error_details", err.Error()))
		return StorageInformation{}, err
	}

	parts := strings.Split(strings.TrimSpace(string(output)), ",")
	if len(parts) < 2 {
		return StorageInformation{Path: driveLetter + ":\\"}, nil
	}

	totalBytes, _ := strconv.ParseUint(strings.TrimSpace(parts[0]), 10, 64)
	freeBytes, _ := strconv.ParseUint(strings.TrimSpace(parts[1]), 10, 64)

	totalMB := totalBytes / (1024 * 1024)
	freeMB := freeBytes / (1024 * 1024)
	usedMB := totalMB - freeMB
	var usedPercent float64
	if totalMB > 0 {
		usedPercent = (float64(usedMB) / float64(totalMB)) * 100
	}

	return StorageInformation{
		Path:        driveLetter + ":\\",
		TotalMB:     totalMB,
		FreeMB:      freeMB,
		UsedPercent: usedPercent,
	}, nil
}