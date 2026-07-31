//go:build linux

package drivers

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

func RetrieveInstalledDrivers() ([]DriverInformation, error) {
	file, err := os.Open("/proc/modules")
	if err != nil {
		slog.Error("Failed to open /proc/modules", slog.String("error_details", err.Error()))
		return nil, fmt.Errorf("failed to read modules: %w", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			slog.Error("Failed to close file", slog.String("error_details", err.Error()))
		}
	}()

	var drivers []DriverInformation
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) > 0 {
			name := fields[0]
			drivers = append(drivers, DriverInformation{
				DeviceName:    name,
				DriverVersion: "",
				Manufacturer:  "Linux Kernel",
				DriverName:    name,
			})
		}
	}
	return drivers, nil
}