//go:build linux

package cpu

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"
)

func RetrieveCPUInfo() (CPUInformation, error) {
	info := CPUInformation{
		Cores: runtime.NumCPU(),
	}

	file, err := os.Open("/proc/cpuinfo")
	if err != nil {
		slog.Error("Failed to open /proc/cpuinfo", slog.String("error_details", err.Error()))
		return info, fmt.Errorf("failed to read cpu info: %w", err)
	}
	
	defer func() {
		if err := file.Close(); err != nil {
			slog.Error("Failed to close file", slog.String("error_details", err.Error()))
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "model name") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				info.ModelName = strings.TrimSpace(parts[1])
				break
			}
		}
	}

	return info, nil
}