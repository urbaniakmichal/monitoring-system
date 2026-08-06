//go:build linux

package memory

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

func (*HardwareMemory)RetrieveMemoryInfo() (MemoryInformation, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		slog.Error("Failed to open /proc/meminfo", slog.String("error_details", err.Error()))
		return MemoryInformation{}, fmt.Errorf("failed to read meminfo: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			slog.Error("Failed to close file", slog.String("error_details", err.Error()))
		}
	}()

	var totalKB, availableKB uint64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		val, _ := strconv.ParseUint(fields[1], 10, 64)
		if strings.HasPrefix(fields[0], "MemTotal:") {
			totalKB = val
		} else if strings.HasPrefix(fields[0], "MemAvailable:") {
			availableKB = val
		}
	}

	totalMB := totalKB / 1024
	availableMB := availableKB / 1024
	usedMB := totalMB - availableMB
	var usedPercent float64
	if totalMB > 0 {
		usedPercent = (float64(usedMB) / float64(totalMB)) * 100
	}

	return MemoryInformation{
		TotalMB:     totalMB,
		AvailableMB: availableMB,
		UsedPercent: usedPercent,
	}, nil
}